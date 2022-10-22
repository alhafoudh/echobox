# echobox [![Build Status](https://img.shields.io/github/workflow/status/alhafoudh/echobox/Docker%20Image%20CI?event=push)](https://github.com/alhafoudh/echobox/actions?query=event%3Apush)

`echobox` is http server that allows you to test http communication inside cloud deployments using customizable http responses.

## Container image

[alhafoudh/eien](https://hub.docker.com/r/alhafoudh/echobox)

## Usage

    ```
    $ docker run -t -i --rm -p 8090:8090 -e TEMPLATE="The foo param: {{ params.foo }}" alhafoudh/echobox
    $ curl http://localhost:8090/?foo=bar
    ```

## License

This project is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).