Welcome to JS Variable Extractor! 🎯
This Go script is designed for developers, security analysts, and code enthusiasts who need to extract and analyze variable declarations and their values from JavaScript files. It efficiently handles multiple URLs, supports both standard input and command-line arguments, and provides options for filtering specific variables.

Features
Concurrent Processing: Leverages a pool of worker threads to process multiple URLs simultaneously, enhancing performance.
Variable Extraction: Uses regex patterns to detect and extract variables and their assigned values from JavaScript files.
Filter Support: Allows you to filter by variable name to focus on specific variables.
Color-Coded Output: Provides clear and visually distinct console output:
Green for detected variables
Yellow for headers and warnings
Red for errors
Usage
From Standard Input
You can pipe URLs into the script's standard input. Optionally, filter results by a specific variable name:
```
[URL] https://example.com/script.js
  Variables and Values:
    [MY_VAR] [12345]
    [ANOTHER_VAR] [someValue]
```
If no variables are found, you will see:
```
[URL] https://example.com/script.js
  No matching variables found.
```
Requirements
Go 1.18 or higher
Contribution
We welcome contributions! Feel free to open issues, submit pull requests, or suggest improvements.
