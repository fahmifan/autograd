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


## Authorization
### policy 1
as a student i can only create, read, edit, & delete my own submissions

### policy 2
as an admin i can only read any submissions from any students

### implementation
```
// auth
polices := [Role][Resource][Action]{
  Admin: {
    User: [ GetAny ]
    Assignment: [ CreateAny, GetAny, GradeAny ],
  },
  Student: {
    User: [Get, Edit],
    Submission: [Create, Edit, Get, Delete],
  },
}

policesV2 := [Role][Permission]{
  Teacher: [ ViewAnyStudents, CreateAnyAssignments, ViewAnyAssignments, ViewAnySubmissions ],
  Student: [ ViewSelfUser, EditSelfUser, CreateSelfSubmission, EditSelfSubmission, ViewSelfSubmission, DeleteSelfSubmission, ViewAnyAssignments ],
}
```

Bgmn best practice untuk handle authorization ?
Saat ini policy di-enforce di tiap handler, krn tiap `resouce` terbatas untuk tiap `roles` nya.
Bisa ga ya enforce policy nya di-centralized, misal di middleware ?

```go
polices := Policies{
  Teacher: [ ViewAnySubmissions, GradeAnySubmissions ],
  Student: [ ViewOwnSubmission ],
}


func handleViewSubmission() {
  if !submission.IsOwnedBy(user) && !authorized(user, ViewAnySubmissions) {
    return unauthorized
  }

  // process ...
}
```
