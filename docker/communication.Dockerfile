FROM default AS builder

WORKDIR /communication

RUN chmod ugo+x .bin/communication

CMD [".bin/communication"]