FROM ubuntu:20.04

ENV PATH /usr/local/bin:$PATH
ENV PATH /usr/local/go/bin:$PATH

RUN apt-get update -y \
    && apt-get install python3 -y \
    && apt-get install python3-dev -y \
    && apt install build-essential -y \
    && apt-get install manpages-dev -y \
    && apt-get install git -y
    # && apt-get install golang-go -y


ADD . /go/src/github.com/facebookresearch/Clinical-Trial-Parser
# RUN cd /go/src/github.com/facebookresearch/Clinical-Trial-Parser && go get -v -t -d ./...

RUN pip3 install -r /go/src/github.com/facebookresearch/Clinical-Trial-Parser/requirements.txt

CMD ["python3"]