# awsping

[![GitHub Actions](https://github.com/yokawasa/awsping/workflows/VS%20Code%20extension%20CI/badge.svg)](https://github.com/yokawasa/awsping/actions)

awsping is a command line tools that reports median latency to
Amazon Web Services regions. It is a fork of [gcping](https://github.com/GoogleCloudPlatform/gcping).

```
awsping [options...]

Options:
-n   Number of requests to be made to each region.
     By default 10; can't be negative.
-c   Max number of requests to be made at any time.
     By default 10; can't be negative or zero.
-t   Timeout. By default, no timeout.
     Examples: "500ms", "1s", "1s500ms".
-top If true, only the top region is printed.

-csv CSV output; disables verbose output.
-v   Verbose output.
```

An example output:

```
$ awsping

 1.  [ap-northeast-1]  50.247016ms
 2.  [ap-northeast-2]  95.119512ms
 3.  [ap-southeast-1]  164.324225ms
 4.  [ap-southeast-2]  234.137157ms
 5.  [us-west-1]       235.008422ms
 6.  [us-west-2]       269.134435ms
 7.  [ap-south-1]      279.4312ms
 8.  [us-east-2]       350.404548ms
 9.  [ca-central-1]    377.65043ms
10.  [eu-west-3]       505.165611ms
11.  [eu-west-2]       515.169295ms
12.  [eu-central-1]    518.303886ms
13.  [eu-west-1]       535.08608ms
14.  [eu-north-1]      575.339622ms
15.  [sa-east-1]       616.53445ms
```

## Installation

* Linux 64-bit: https://awspingrelease.s3-ap-northeast-1.amazonaws.com/awsping_linux_amd64
  ```
  $ curl https://awspingrelease.s3-ap-northeast-1.amazonaws.com/awsping_linux_amd64 > awsping && chmod +x awsping
  ```
* Mac 64-bit: https://awspingrelease.s3-ap-northeast-1.amazonaws.com/awsping_darwin_amd64
  ```
  $ curl https://awspingrelease.s3-ap-northeast-1.amazonaws.com/awsping_darwin_amd64 > awsping && chmod +x awsping
  ```

* Windows 64-bit: https://awspingrelease.s3-ap-northeast-1.amazonaws.com/awsping_windows_amd64

Or, you can always build the binary from the source code like this:

```
$ git clone https://github.com/yokawasa/awsping.git
$ cd awsping
$ make
$ tree bin

bin
├── awsping_darwin_amd64
├── awsping_linux_amd64
└── awsping_windows_amd64
```

## GitHub Actions Release Workflow

By your pushing tag, GitHub trigger the GitHub Actions Release workflow where
- The project is checkout and build in multi-SO & Architecture
- Release each artifact to release URL in the repository

This is how you trigger the workflow
```
git tag -a v0.0.1 -m "Version awsping-v0.0.1"
git push --tags
```

See [the workflow]() for the detail
