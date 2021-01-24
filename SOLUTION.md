# Bootstrapping and running

Solution is built and packed using docker, therefore please make sure
[docker is installed](https://docs.docker.com/get-docker/) on your host.

## Easy way

Most easy way to run the app is using `run.sh` which will build an image and run the container with mounted project
directory.

It is expected that input file is a plain json. It could be either file located in the local filesystem or valid URL to
download the file.

If you want to mount another folders into the container you can specify those in `VOLUMES` array in `run.sh`. Should you
provide input source from different mounted folder, make sure specify absolute file path.

You can control over search params via input arguments, please find full list below:

```shell
bash-3.2$ ./run.sh
Source json file is not provided, please see usage for available options:
Usage of collector:
  -postcode string
        Search for a given postcode
  -recipe value
        Search for a recipes containing given pattern, can be specified multiple times
  -src string
        Source json file path with deliveries data to collect
  -window-from string
        Search for a given postcode in window from (default "12AM")
  -window-to string
        Search for a given postcode in window to (default "11PM")
```

### Input from local file

Example of running recipe collector with all the possible arguments. We want to run it against demo file which is
located in Downloads folder
(i.e. folder outside the project root):

```shell
./run.sh run --src /Users/username/Downloads/hf_test_calculation_fixtures.json --recipe steak --recipe honey --postcode 10163 --window-from 10AM --window-to 11PM
```

To make sure filepath `/Users/username/Downloads/hf_test_calculation_fixtures.json`
is reachable from within the container I'll update `run.sh` as follows:

```shell
VOLUMES=(
  "${WD}"
  "/Users/username/Downloads" # This volume is added, and will be mapped with the exact same path into container
)
```

### Input from URL

```shell
./run.sh run --src https://pastebin.com/raw/6gtp4Qyc --recipe steak --recipe honey --recipe tex --postcode 10163 --window-from 10AM --window-to 11P
```

## Harder way

You can easily build image and run container on your own.

Dockerfile provides 2 targets: **build** and **runtime** which are based on different base images. If you need go
installed within container (e.g. to run tests), use **build**
target, otherwise **runtime** will give you option to run compiled app within a minimalistic alpine environment.

In both, **build** and **runtime** targets you can find app in global PATH as `recipe-collector`.

# Known issues

### MacOS `zsh`

Running container in MacOS zsh will produce extra `%` symbol in the end, which makes output json invalid.

You can easily spawn `bash` subshell and run container within, all should be fine this way.

### Alpine repositories issue

Sometimes during image build you may face an error related to alpine repositories update:

```shell
ERROR: https://dl-cdn.alpinelinux.org/alpine/v3.13/main: temporary error (try again later)
```

Just run the same command once again, and issue should disappear

# Implementation details

Solution is quite simplified, which means certain mandatory for production must-have features are omitted, e.g. some
input validation, proper error handling, logging, unexpected data processing. Main goal was to fulfil main functional
requirements.

I tried to focus on collectors design with performance not being affected too much, all besides that was skipped in sake
of time saving. Several unit tests are in place for demo purposes as well.

# Possible improvements

Standard encode/json lib might be not the optimal solution in terms of performance for json stream parsing due to struct
mapping and reflection used. Unfortunately I'm not really familiar with third-party libs that might provide same stream
parsing capabilities without reflection overhead. Quick googling didn't help either to find solution that I could apply
as quick as encode/json.

Validation of input arguments is barely done now, and it's not even an improvement suggestion but must have requirement
for production environment. Same applies for exceptional cases handling and logging, e.g. unexpected data structures,
unexpected values format, out of expected range values, etc. 