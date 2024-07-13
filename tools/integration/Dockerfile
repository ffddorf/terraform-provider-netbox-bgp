ARG NETBOX_VERSION
FROM netboxcommunity/netbox:v${NETBOX_VERSION}

ARG NETBOX_BGP_VERSION
RUN /opt/netbox/venv/bin/pip install --no-warn-script-location netbox-bgp==${NETBOX_BGP_VERSION}

COPY plugins.py /etc/netbox/config/plugins.py
RUN SECRET_KEY="dummydummydummydummydummydummydummydummydummydummy" /opt/netbox/venv/bin/python /opt/netbox/netbox/manage.py collectstatic --no-input
