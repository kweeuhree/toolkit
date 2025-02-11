# 🛠️ Toolkit Package

## 🔍 Overview

This repository contains a toolkit with reusable functions for common tasks encountered in web development with Go.
The goal is to avoid external dependencies and provide thorough test coverage.

## 📄 Description

The toolkit package provides essential helper methods for generating random strings and handling file uploads. The key functionality includes:

- **Random String Generation** generates cryptographically secure random strings for various use cases.
- **File Uploads** upload of one or more files, check file type validity, and save them to a directory with optional renaming.
- **JSON functionalilty** provides essential helper methods for handling JSON in HTTP requests and responses, such as read and write JSON files, send a JSON error response.

**Test coverage**: 89.9% of statements.

## Installation

```bash
go get -u github.com/kweeuhree/toolbox
```

## ⚙️ Functions

#### ➡️RandomString

Generates a random string of length n using a predefined set of characters (randomStrSource). The string is generated cryptographically securely using Go's crypto/rand package.

**Parameters**:

- `n`: The length of the random string to be generated.

**Returns**:

- A random string of length n.

**Example**:

```go
t := &toolkit.Tools{}
randomStr := t.RandomString(16)
fmt.Println(randomStr)  // e.g., "aB3fGz0a8sK2LsD8"
```

#### ➡️UploadFiles

Uploads one or more files to a specified directory and gives the files a random name.

**Parameters**:

- `r`: The HTTP request containing the files to upload.
- `uploadDir`: The directory where the files should be uploaded.
- `rename`: (Optional) If set to false, the files will keep their original names. If not provided, files will be renamed using random strings.

**Returns**:

- A slice of `UploadedFile` structs containing details about the uploaded files.
- An error, if any occurred during the upload.

**Example**:

```go
uploadedFiles, err := t.UploadFiles(r, "./uploads")
if err != nil {
  fmt.Println("Error uploading files:", err)
}
for _, file := range uploadedFiles {
  fmt.Println("Uploaded file:", file.NewFileName)
}
```

#### ➡️UploadOneFile

Convenience method that calls `UploadFiles` expecting only one file to be uploaded.

**Parameters**:

- `r`: The HTTP request containing the file to upload.
- `uploadDir`: The directory where the file should be uploaded.
- `rename`: (Optional) If set to false, the file will keep its original name.

**Returns**:

- An `UploadedFile` containing details about the uploaded file.
- An error, if any occurred during the upload.

**Example**:

```go
uploadedFile, err := t.UploadOneFile(r, "./uploads")
if err != nil {
    fmt.Println("Error uploading file:", err)
}
fmt.Println("Uploaded file:", uploadedFile.NewFileName)
```

#### ➡️ ReadJSON

Reads and decodes JSON data from an HTTP request body into the provided 'data' object. It validates the JSON format, checks the request size, and handles various error scenarios, including syntax errors, unknown fields, and unexpected EOF.

**Parameters**:

- `w`: The HTTP response writer.
- `r`: The HTTP request containing the JSON body.
- `data`: A pointer to the struct where the decoded JSON data will be stored.

**Returns**:

- An error if the JSON is malformed or the body exceeds the allowed size.

**Example**:

```go
t := &toolkit.Tools{}
var myData MyStruct
err := t.ReadJSON(w, r, &myData)
if err != nil {
    fmt.Println("Error:", err)
}
```

#### ➡️ WriteJSON

Writes a JSON response with the provided status, data, and optional custom headers.

**Parameters**:

- `w`: The HTTP response writer.
- `status`: The HTTP status code for the response.
- `data`: The data to be written as JSON.
- `headers`: Optional custom headers to include in the response.

**Returns**:

- An error if the response writing fails.

**Example**:

```go
t := &toolkit.Tools{}
responseData := JSONResponse{Error: false, Message: "Success"}
err := t.WriteJSON(w, http.StatusOK, responseData, customHeaders)
if err != nil {
    fmt.Println("Error:", err)
}
```

#### ➡️ ErrorJSON

Sends a JSON error response with an optional custom status code.

**Parameters**:

- `w`: The HTTP response writer.
- `err`: The error to be included in the response.
- `status`: Optional HTTP status code (default is 400).

**Returns**:

- An error if writing the error response fails.

**Example**:

```go
t := &toolkit.Tools{}
err := t.ErrorJSON(w, fmt.Errorf("something went wrong"))
if err != nil {
    fmt.Println("Error:", err)
}
```

## 🚩 Error Handling

UploadFiles and UploadOneFile will return an error if:

- The file type is not allowed (checked against AllowedFileTypes).
- The file size exceeds the configured MaxFileSize.
- There are issues opening or saving the file.
  Make sure to handle these errors appropriately in your application.

## 📦 Dependencies

This toolkit relies only on standard Go packages.

## ✅ Testing

The toolkit uses built-in `testing` package for testing app logic, and provides comprehensive tests that ensure that the functions behave correctly for different input conditions and handle edge cases, such as file type restrictions and renaming.

Coverage: 89.9% of statements

## 📃 License

This package is licensed under the MIT License. See LICENSE for more information.
