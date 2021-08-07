# Autograde

This WIP is a autograd backend
Autograd is a CLI that run a c++ source code from `submission` and test it againts the `input` & `ouput`

## Notes

- the submission name must follow this: `{userID}`-`{test_code}`-.cpp
- `test-code` is file name used for input & output. Use underscore to spereate between words

## How to

- compatible os: *nix family should do
  - tested on: Mac OS Catalina version 10.15.4
- you need `golang v1.14` & `g++` installed in your machine
- to run the code: `make grade`
