FROM postgres:14.6-alpine3.16

ENV PERL_MM_USE_DEFAULT=1
ENV PG_CRON_VERSION=1.4.0

RUN apk add --no-cache --update \
        build-base \
        curl \
        expat \
        expat-dev \
        git \
        gnupg \
        heirloom-mailx \
        llvm \
        llvm-dev \
        make \
        perl \
        perl-app-cpanminus \
        perl-dev \
        perl-dbd-pg \
        perl-dbi \
        perl-xml-libxml \
        perl-xml-parser \
        perl-xml-simple \
        postgresql-client \
        postgresql-dev \
        sudo \
        wget

RUN     rm -rf /var/cache/apk/* /tmp/*

RUN     cpan -i Module::Signature \
                TAP::Parser::SourceHandler::pgTAP \
                XML::SAX::Expat \
# install pg_prove
RUN cpanm TAP::Parser::SourceHandler::pgTAP
RUN cpanm Test::Deep
RUN cpanm TAP::Harness::JUnit

WORKDIR /var/lib/postgresql

# install pgtap
RUN git clone https://github.com/theory/pgtap.git \
    && cd pgtap \
    && make \
    && make install \
    && make clean
# install pg_partman
RUN git clone https://github.com/pgpartman/pg_partman.git \
    && cd pg_partman \
    && make NO_BGW=1 install

# install pg_cron
RUN wget -O /var/lib/postgresql/pg_cron.tgz https://github.com/citusdata/pg_cron/archive/v$PG_CRON_VERSION.tar.gz \
    && tar xvzf /var/lib/postgresql/pg_cron.tgz \
    && cd pg_cron-$PG_CRON_VERSION \
    && sed -i.bak -e 's/-Werror//g' Makefile \
    && sed -i.bak -e 's/-Wno-implicit-fallthrough//g' Makefile \
    && export PATH=/usr/bin:$PATH \
    && make \
    && sudo PATH=$PATH make install \
    && cp /usr/lib/postgresql14/pg_cron.so /usr/local/lib/postgresql/ \
    && cp /usr/share/postgresql14/extension/* /usr/local/share/postgresql/extension

CMD ["-c", "shared_preload_libraries=pg_cron", "-c", "cron.database_name=scam"]