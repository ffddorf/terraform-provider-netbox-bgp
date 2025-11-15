ARG NETBOX_VERSION
FROM netboxcommunity/netbox:v${NETBOX_VERSION}

ARG NETBOX_BGP_VERSION
RUN /usr/local/bin/uv pip install netbox-bgp==${NETBOX_BGP_VERSION}

COPY plugins.py /etc/netbox/config/plugins.py
RUN SECRET_KEY="dummydummydummydummydummydummydummydummydummydummy" /opt/netbox/venv/bin/python /opt/netbox/netbox/manage.py collectstatic --no-input
