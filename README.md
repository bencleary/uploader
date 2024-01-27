# File Upload Service

## Overview

A small service which can handle image uploads for an imaginary chat application. Based on the implementation discussed by [Code Aesthetic](https://www.youtube.com/watch?v=J1f5b4vcxCQ). In this video implementation is done leveraging TypeScript, as a learning exercise i wanted to implement it using another language and chose Go.

### Assumptions

- Only one storage engine is needed as the service runs, unlike the example which routes based on client ID.
- The implementation assumes that user authentication has already been taken care of by another service.
- JSON will be used as the response type.

## Features

- File encryption using AES ✅
- Image resizing ✅
- Preview generation ✅
- Filing for historical access ✅

## Structure

All files wituin the root directory are types used by the service in other layers of the application. Business logic is split throughout the layers of the application and is stiched to gether in the HTTP handler upload, and download. The internal directory holds the concreate implementations for these layers such as DB, http, etc.

The approach the application takes is to leverage dependancy inversion, dependency injection and the startegy pattern to make a flexible and extensible service.

## CLI

Provides a CLI for uploading files on a adhoc basis

## HTTP

Provides a rest api which enables files to be uploaded directly.

### Endpoints

#### Upload

`/file/upload`

#### ID

`/file/:id`

## Usage

- Post method to file/upload
- Request body is made up of file and key
- key must be 16 bit
- 

## Tasks

- [ ] Better error messages in packages - 1
- [ ] Genericise Filer - 2
- [ ] Bring in user specific keys not encryption specific keys - 9
- [x] Key should be in header
- [ ] Add tests - 4
- [ ] Add S3 storage - 5
- [ ] Add CLI interface - 6
- [ ] Godocs for methods - 8
- [ ] Replace ReadCloser with more generic interface - 3
- [ ] Implement queue for processing files - 7


