# Performance Testing Guide

This document explains how to test the performance of the todo application with large JSON files.

## Overview

The application currently loads entire JSON files into memory, which can impact performance with large datasets. This guide helps you:

1. **Benchmark** JSON loading, saving, and sorting operations
2. **Generate** test files of various sizes
3. **Identify** performance bottlenecks
4. **Understand** scaling characteristics

---

## Quick Start

```bash
# Run all benchmarks
go test -bench=. -benchmem ./todo

# Run specific benchmark
go test -bench=BenchmarkLoad_Large -benchmem ./todo

# Run performance tests
go test -v ./todo -run TestLargeFile

# Generate test data
go run scripts/generate_test_data.go -count 10000
```

---

## Test Files Structure

### Benchmark Tests (`todo/benchmark_test.go`)

Measures **raw performance** of operations:

- **Load benchmarks**: `BenchmarkLoad_*` - JSON unmarshaling from disk
- **Save benchmarks**: `BenchmarkSave_*` - JSON marshaling to disk (atomic writes)
- **Sort benchmarks**: `BenchmarkSort_*` - Sorting incomplete/completed todos
- **CRUD benchmarks**: `BenchmarkAddTodo_*`, `BenchmarkToggleTodo_*`
- **JSON benchmarks**: `BenchmarkJSONMarshal_*`, `BenchmarkJSONUnmarshal_*`

**Sizes tested**: Small (100), Medium (1,000), Large (10,000), VeryLarge (100,000)

### Performance Tests (`todo/performance_test.go`)

Validates **functional correctness** and **performance thresholds**:

- `TestLargeFileLoad` - Load time should be < 2 seconds
- `TestLargeFileSave` - Save time should be < 3 seconds
- `TestLargeFileIntegrity` - Data integrity after save/load cycle
- `TestMemoryUsage` - Memory consumption estimates
- `TestSortPerformance` - Sort time should be < 100ms
- `TestFileCorruption` - Corrupted file recovery
- `TestJSONFileSize` - File size analysis

---

## Running Tests

### 1. Benchmark Tests

```bash
# Run all benchmarks with memory allocation stats
go test -bench=. -benchmem ./todo

# Run only Load benchmarks
go test -bench=BenchmarkLoad -benchmem ./todo

# Run with CPU profiling
go test -bench=BenchmarkLoad_Large -benchmem -cpuprofile=cpu.prof ./todo

# Run with memory profiling
go test -bench=BenchmarkLoad_Large -benchmem -memprofile=mem.prof ./todo

# View profiling results
go tool pprof cpu.prof
# In pprof: top, list, web
```

### 2. Performance Tests

```bash
# Run all performance tests
go test -v ./todo -run TestLarge

# Run specific test
go test -v ./todo -run TestLargeFileLoad

# Skip memory tests (use -short flag)
go test -v -short ./todo
```

### 3. Generate Test Data

```bash
# Generate 1,000 todos (default)
go run scripts/generate_test_data.go

# Generate 50,000 todos
go run scripts/generate_test_data.go -count 50000

# Custom output path
go run scripts/generate_test_data.go -count 10000 -output /tmp/test.json

# 50% completion rate (default is 33%)
go run scripts/generate_test_data.go -count 5000 -completion 50
```

---

## Understanding Benchmark Results

### Example Output

```
BenchmarkLoad_Small-8            5000    245123 ns/op    125000 B/op    1500 allocs/op
BenchmarkLoad_Medium-8            500   2451234 ns/op   1250000 B/op   15000 allocs/op
BenchmarkLoad_Large-8              50  24512345 ns/op  12500000 B/op  150000 allocs/op
```

**Columns:**
- `BenchmarkLoad_Small-8`: Test name with GOMAXPROCS
- `5000`: Number of iterations run
- `245123 ns/op`: Nanoseconds per operation (245.123 µs)
- `125000 B/op`: Bytes allocated per operation (122 KB)
- `1500 allocs/op`: Number of allocations per operation

### Performance Metrics

| Size | Todos | Expected Load Time | Expected File Size | Memory Usage |
|------|-------|--------------------|--------------------|--------------|
| Small | 100 | < 1 ms | ~10 KB | ~10 KB |
| Medium | 1,000 | < 10 ms | ~100 KB | ~100 KB |
| Large | 10,000 | < 100 ms | ~1 MB | ~1 MB |
| VeryLarge | 100,000 | < 1 s | ~10 MB | ~10 MB |

---

## Performance Bottlenecks

### Current Implementation Issues

1. **Full File Loading** (`todo/todo.go:150-170`)
   - Entire JSON file loaded into memory via `json.Unmarshal`
   - **Impact**: O(n) memory usage, slow startup with large files
   - **Location**: `todo.Load()`

2. **No Pagination** (`ui/view.go:159-206`)
   - All todos rendered in loop
   - **Impact**: O(n) rendering, UI lag with 10,000+ items
   - **Location**: `renderTodoList()`

3. **Linear Sort** (`todo/todo.go:92-108`)
   - Manual partition into incomplete/completed
   - **Impact**: O(n) on every toggle/add operation
   - **Location**: `TodoList.Sort()`

4. **Atomic Writes** (`todo/todo.go:111-147`)
   - Temp file + rename for safety
   - **Impact**: Disk I/O overhead, but necessary for data integrity
   - **Location**: `TodoList.Save()`

### Scaling Characteristics

| Operation | Time Complexity | Space Complexity | Bottleneck |
|-----------|-----------------|------------------|------------|
| Load | O(n) | O(n) | JSON parsing |
| Save | O(n) | O(n) | JSON marshaling + disk I/O |
| Sort | O(n) | O(1) | Linear scan |
| Render | O(n) | O(1) | Terminal rendering |
| Add/Toggle | O(1) + O(n) sort | O(1) | Sort after mutation |

---

## Optimization Strategies

### Short-term (Easy Wins)

1. **Lazy Rendering**: Only render visible todos in viewport
   ```go
   // Instead of rendering all todos, render viewport only
   start := m.scrollOffset
   end := min(start + m.viewportHeight, len(todos))
   for i := start; i < end; i++ { ... }
   ```

2. **Binary Search for Completed Boundary**
   ```go
   // Instead of linear scan in Sort(), track boundary index
   type TodoList struct {
       Todos []Todo
       CompletedIndex int // First completed todo index
   }
   ```

3. **Incremental Sort**: Only re-sort on toggle, not on render

### Mid-term (Moderate Effort)

1. **Streaming JSON Parser**: Use `json.Decoder` for large files
2. **Virtual Scrolling**: Only keep visible items in memory
3. **Caching**: Cache rendered strings, invalidate on change
4. **Background Loading**: Load files asynchronously

### Long-term (Major Refactor)

1. **Database Backend**: SQLite instead of JSON files
2. **Pagination**: Limit loaded todos to 100-500 at a time
3. **Indexing**: Add search indices for fast lookup
4. **Compression**: Gzip JSON files for storage

---

## Expected Performance Impact

### At Different Scales

**100 todos (typical personal use):**
- ✅ No performance issues
- Load: < 1ms, Render: instant, Sort: negligible

**1,000 todos (power user):**
- ✅ Minor impact
- Load: ~10ms, Render: slight delay, Sort: ~1ms

**10,000 todos (stress test):**
- ⚠️ Noticeable lag
- Load: ~100ms, Render: visible delay, Sort: ~10ms
- **Action needed**: Implement lazy rendering

**100,000 todos (extreme):**
- ❌ Unusable
- Load: ~1s, Render: seconds, Sort: ~100ms
- **Action needed**: Major refactor (database, pagination)

### Real-world Usage

Most users will have **50-500 todos** across multiple files:
- Current implementation is **sufficient** for typical use
- Optimizations beneficial for power users (1,000+)
- Database migration only needed if targeting 10,000+ todos

---

## Continuous Monitoring

### Add to CI/CD

```bash
# In GitHub Actions or similar
- name: Run benchmarks
  run: |
    go test -bench=. -benchmem ./todo > benchmark_results.txt

# Track regression over time
- name: Compare benchmarks
  run: |
    # Compare against baseline
    benchstat baseline.txt benchmark_results.txt
```

### Benchstat Tool

```bash
# Install benchstat
go install golang.org/x/perf/cmd/benchstat@latest

# Run benchmarks twice
go test -bench=. -count=5 ./todo > old.txt
# Make changes...
go test -bench=. -count=5 ./todo > new.txt

# Compare
benchstat old.txt new.txt
```

---

## Profiling Deep Dive

### CPU Profile

```bash
go test -bench=BenchmarkLoad_Large -cpuprofile=cpu.prof ./todo
go tool pprof -http=:8080 cpu.prof
```

### Memory Profile

```bash
go test -bench=BenchmarkLoad_Large -memprofile=mem.prof ./todo
go tool pprof -http=:8080 mem.prof
```

### Trace Analysis

```bash
go test -bench=BenchmarkLoad_Large -trace=trace.out ./todo
go tool trace trace.out
```

---

## Questions to Answer

Use these tests to answer:

1. **"How many todos can the app handle?"**
   - Run `TestLargeFileLoad` with increasing sizes
   - Find threshold where load time exceeds 2 seconds

2. **"Which operation is slowest?"**
   - Compare `BenchmarkLoad_Large` vs `BenchmarkSave_Large` vs `BenchmarkSort_Large`
   - Identify bottleneck operation

3. **"How much memory does a 10,000 todo file use?"**
   - Run `TestMemoryUsage` and check logs
   - ~1MB in memory, ~1MB on disk

4. **"Is JSON parsing or disk I/O the bottleneck?"**
   - Compare `BenchmarkJSONUnmarshal_Large` (pure parsing) vs `BenchmarkLoad_Large` (disk + parsing)
   - If similar, parsing is bottleneck; if different, disk I/O is

---

## Summary

- ✅ **Current state**: Handles typical usage (100-1,000 todos) well
- ⚠️ **Scaling limit**: Degrades noticeably at 10,000+ todos
- 🔧 **Low-hanging fruit**: Lazy rendering, incremental sort
- 🚀 **Future-proofing**: Database migration for 50,000+ todos

**Recommended action**: Run benchmarks to establish baseline, then optimize rendering for 10,000+ todo use cases if needed.
