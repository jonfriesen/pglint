FROM golang:1.11

LABEL "name"="golang tools action"
LABEL "maintainer"="Jon Friesen <jon@jonfriesen.ca>"
LABEL "version"="0.0.1"

LABEL "com.github.actions.icon"="code"
LABEL "com.github.actions.color"="green-dark"
LABEL "com.github.actions.name"="Go Tools"
LABEL "com.github.actions.description"="This is an Action to run go commands."

RUN apt-get update
RUN apt-get install -y libecpg-dev

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]