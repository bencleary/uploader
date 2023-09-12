# 🚧 WIP 🚧 - File Upload Service

## Overview

A small service which can handle image uploads for an imaginary service or chat application. Based on the implementation discussed by [Code Aesthetic](https://www.youtube.com/watch?v=J1f5b4vcxCQ). In this video they detail implementation using TypeScript, as a learning exercise i wanted to implement it using another language and chose Go.

### Assumptions

- Only one storage engine is needed as the service runs, unlike the example which routes based on client ID.
- The implementation assumes that user authentication has already been taken care of by another service.
- JSON responses are used, but i would like to extend to using HTMX and serve both JSON and HTML.

## Features

- File encryption using AES ✅
- Image resizing ✅
- Preview generation ✅
- Filing for historical access ✅

## Structure

All files wituin the root directory are types used by the service in other layers of the application. Business logic is split throughout the layers of the application and is stiched to gether in the HTTP handler upload, and download. The internal directory holds the concreate implementations for these layers such as DB, http, etc.

The approach the application takes is to leverage dependancy inversion, dependency injection and the startegy pattern to make a flexible and extensible service.

## Todo

- [ ] Add tests
- [ ] Add configuration
- [ ] Add CLI interface


