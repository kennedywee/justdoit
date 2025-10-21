#!/bin/bash
# Performance testing automation script for justdoit

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Create results directory
RESULTS_DIR="perf_results_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$RESULTS_DIR"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  JustDoIt Performance Test Suite${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Function to print section headers
print_header() {
    echo -e "\n${YELLOW}>>> $1${NC}\n"
}

# Function to print success
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# Function to print error
print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# 1. Run benchmark tests
print_header "Running Benchmark Tests"
go test -bench=. -benchmem ./todo | tee "$RESULTS_DIR/benchmarks.txt"
print_success "Benchmarks saved to $RESULTS_DIR/benchmarks.txt"

# 2. Run performance tests
print_header "Running Performance Tests"
go test -v ./todo -run TestLarge 2>&1 | tee "$RESULTS_DIR/performance_tests.txt"
print_success "Performance tests saved to $RESULTS_DIR/performance_tests.txt"

# 3. Run memory tests (if not in short mode)
print_header "Running Memory Usage Tests"
go test -v ./todo -run TestMemoryUsage 2>&1 | tee "$RESULTS_DIR/memory_tests.txt"
print_success "Memory tests saved to $RESULTS_DIR/memory_tests.txt"

# 4. Generate CPU profile
print_header "Generating CPU Profile"
go test -bench=BenchmarkLoad_Large -cpuprofile="$RESULTS_DIR/cpu.prof" ./todo > /dev/null 2>&1
print_success "CPU profile saved to $RESULTS_DIR/cpu.prof"
echo "  View with: go tool pprof -http=:8080 $RESULTS_DIR/cpu.prof"

# 5. Generate memory profile
print_header "Generating Memory Profile"
go test -bench=BenchmarkLoad_Large -memprofile="$RESULTS_DIR/mem.prof" ./todo > /dev/null 2>&1
print_success "Memory profile saved to $RESULTS_DIR/mem.prof"
echo "  View with: go tool pprof -http=:8080 $RESULTS_DIR/mem.prof"

# 6. Generate test data files
print_header "Generating Test Data Files"
echo "Creating test files in ~/.tui_todos/..."

for size in 100 1000 10000 50000; do
    go run scripts/generate_test_data.go -count $size -output "$HOME/.tui_todos/test_${size}.json" 2>&1 | grep "✓"
done

print_success "Test data files created in ~/.tui_todos/"

# 7. Create summary report
print_header "Creating Summary Report"

SUMMARY_FILE="$RESULTS_DIR/SUMMARY.md"

cat > "$SUMMARY_FILE" <<EOF
# Performance Test Results

**Date:** $(date '+%Y-%m-%d %H:%M:%S')
**Go Version:** $(go version)
**OS:** $(uname -s) $(uname -r)

---

## Benchmark Results

\`\`\`
$(cat "$RESULTS_DIR/benchmarks.txt" | grep "Benchmark")
\`\`\`

---

## Performance Test Results

\`\`\`
$(grep -E "(PASS|FAIL|---)" "$RESULTS_DIR/performance_tests.txt" | head -20)
\`\`\`

---

## Key Metrics

### Load Times

$(grep "Loaded.*todos in" "$RESULTS_DIR/performance_tests.txt" || echo "No load time data found")

### Save Times

$(grep "Saved.*todos" "$RESULTS_DIR/performance_tests.txt" || echo "No save time data found")

### Sort Times

$(grep "Sorted.*todos in" "$RESULTS_DIR/performance_tests.txt" || echo "No sort time data found")

### File Sizes

$(grep "bytes/todo" "$RESULTS_DIR/performance_tests.txt" || echo "No file size data found")

---

## Memory Usage

\`\`\`
$(grep -A2 "Loaded.*todos - File size" "$RESULTS_DIR/memory_tests.txt" || echo "No memory data found")
\`\`\`

---

## Profiles

- **CPU Profile:** $RESULTS_DIR/cpu.prof
- **Memory Profile:** $RESULTS_DIR/mem.prof

View with:
\`\`\`bash
go tool pprof -http=:8080 $RESULTS_DIR/cpu.prof
go tool pprof -http=:8080 $RESULTS_DIR/mem.prof
\`\`\`

---

## Test Data Files

Generated test files in \`~/.tui_todos/\`:
- test_100.json
- test_1000.json
- test_10000.json
- test_50000.json

---

## Recommendations

Based on benchmark results:

1. **Current Performance:**
   - Review load times in the benchmark section above
   - Check if any operation exceeds thresholds in PERFORMANCE_TESTING.md

2. **Optimization Priorities:**
   - If load times > 1s for 10K todos: Implement lazy loading
   - If render times lag: Implement virtual scrolling
   - If sort times > 100ms: Use incremental sorting

3. **Next Steps:**
   - Profile with: \`go tool pprof -http=:8080 $RESULTS_DIR/cpu.prof\`
   - Identify hotspots in flame graph
   - Compare against future runs to track regressions

EOF

print_success "Summary report saved to $RESULTS_DIR/SUMMARY.md"

# 8. Display summary
print_header "Test Summary"

echo -e "${BLUE}Results Directory:${NC} $RESULTS_DIR"
echo ""
echo -e "${GREEN}Files created:${NC}"
ls -lh "$RESULTS_DIR"
echo ""

# Quick stats
BENCHMARK_COUNT=$(grep -c "^Benchmark" "$RESULTS_DIR/benchmarks.txt" || echo "0")
TEST_PASS=$(grep -c "PASS:" "$RESULTS_DIR/performance_tests.txt" || echo "0")
TEST_FAIL=$(grep -c "FAIL:" "$RESULTS_DIR/performance_tests.txt" || echo "0")

echo -e "${BLUE}Statistics:${NC}"
echo "  Benchmarks run: $BENCHMARK_COUNT"
echo "  Tests passed: $TEST_PASS"
echo "  Tests failed: $TEST_FAIL"
echo ""

if [ $TEST_FAIL -gt 0 ]; then
    print_error "Some tests failed! Check $RESULTS_DIR/performance_tests.txt for details"
else
    print_success "All tests passed!"
fi

echo ""
echo -e "${YELLOW}View full report:${NC} cat $RESULTS_DIR/SUMMARY.md"
echo -e "${YELLOW}View benchmarks:${NC} cat $RESULTS_DIR/benchmarks.txt"
echo -e "${YELLOW}Profile CPU:${NC} go tool pprof -http=:8080 $RESULTS_DIR/cpu.prof"
echo -e "${YELLOW}Profile Memory:${NC} go tool pprof -http=:8080 $RESULTS_DIR/mem.prof"
echo ""
echo -e "${GREEN}Done!${NC}"
