# BENCOMP
bencomp (BENchmark COMPression) is a free, open-source command-line utility for estimating and comparing the performance of various compression algorithms. It is intended to help quickly determine what compression algorithm may be best suited for your application. Bear in mind that this tool merely provides an estimate; there is no substitute for proper profiling and benchmarking in a real environment.
```
> ./bencomp.exe --rand-gen --json-str-len 100 --json-num-fields 50 --json-max-depth 5
Original data size: 14.0817 MB
Compression-Library      Total-Time    Compressed-Size    Ratio
gzip                     361.4089ms    8.3811 MB          59.52%
zlib-default             359.6119ms    8.3811 MB          59.52%
zlib-best-compression    355.7197ms    8.3811 MB          59.52%
zlib-best-speed          145.5566ms    8.4687 MB          60.14%
zstd                     67.8633ms     10.3429 MB         73.45%
```

## Installation
bencomp is available at `https://github.com/cameron-p-wilson/bencomp`. You can either:
1. Clone the source code, then build the executable using `go build`.
2. Download the executable for your OS directly from GitHub.

## Getting Started
To run a benchmark, you must supply some input. You can give either a file for bencomp to read, or you can tell bencomp to randomly generate data for you.
You can use `bencomp -h` to get a list of all options.

- `bencomp --rand-gen` -- this will direct bencomp to randomly generate JSON data to use in compression and decompression.
- `bencomp --file` -- this will direct bencomp to read data from the given file to use in compression and decompression

### JSON Generation
If you have an idea of the kind of JSON payloads that your application is likely to deal with, you can direct bencomp to randomly generate JSON in a similar pattern. All flags which affect JSON generation have the `json-` prefix.

The JSON generator will only create string values with lowercase characters. The JSON tree can have any number of fields and any level of nesting as long as the JSON does not exceed `2^32` bytes in size.

#### Examples
Each example below describes a way to construct a JSON tree of varying patterns; the JSON is randomly generated then immediately used in benchmarking.
 - `bencomp --rand-gen --json-num-fields 6 --json-degree 3 --json-max-depth 4`
    - Creates a JSON tree where each node has exactly 6 key-value pairs and 3 children, and the depth of any branch in the tree shall not exceed 4. (Be careful with large numbers of `--json-max-depth` as the tree grows quickly in size!)
 - `bencomp --rand-gen --json-num-fields-range 0-5 --json-degree-range 0-3 --json-max-depth 8`
    - Creates a JSON tree where each node has between 0 to 5 key-value pairs and between 0 to 3 children, and the depth of any branch in the tree shall not exceed 8. Values are picked randomly in a uniform distribution.
 - `bencomp --rand-gen --json-str-len 64`
    - Creates a JSON tree with default values for `--json-num-fields`, `--json-degree`, and `--json-max-depth`. However, each key and each value will be exactly 64 characters in size.
 - `bencomp --rand-gen --json-str-len-range 6-32`
    - Similar to the above, but each key and each value will be between 6 to 32 characters in size.
 - `bencomp --rand-gen --json-dict-size 10 --json-str-len 32`
    - Creates a JSON tree with the default structure, but each key and each value will be chosen from a set of 10 randomly generated strings, each 32 characters in size. This will result in a significantly smaller compressed file size.
 - `bencomp --rand-gen --json-dict-file <file.txt>`
    - Creates a JSON tree with the default structure, but each key and each value will be chosen from a file. The input file should be a plaintext file containing a separate word on each line and nothing else.

### Optional Statistics
By default, bencomp will display the total time, uncompressed file size, compressed file size, and compression ratio for each compression library used in the benchmark. There are, however, additional options:
 - `bencomp --show-json`
    - Display the data used for compression testing (this can be very large in many cases). If used with `--count` and `--rand-gen`, will display only the last input.
 - `bencomp --show-compress-time`
    - Display the time spent on compression.
 - `bencomp --show-decompress-time`
    - Display the time spent on decompression.
 - `bencomp --network-bandwidth <bandwidth>`
    - Display the time it would take to compress, send over a network, and decompress a certain number of payloads. The bandwidth is expressed in bytes per second, e.g. `1000` = `1000 bytes per second`, or `128MB` = `128000000 bytes per second`. For simplicity, it is assumed that there is only 1 producer, only 1 consumer, and only a single network path with no latency or dropped packets. (This is where you would want to test in a real environment).
 - `bencomp --network-payloads <n>`
    - Used by `--network-bandwidth` in calculating the total time.
 - `bencomp --count <n>`
    - Repeat the benchmark `n` times and report the median values. If used with `--rand-gen`, a new random JSON payload will be generated each time.