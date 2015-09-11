# cf-events

A data store for Cloud.gov / Cloud Foundry [events logs](http://apidocs.cloudfoundry.org/215/events/list_all_events.html). 

### Accessing Data - [cf-ssh](https://docs.18f.gov/getting-started/cf-ssh/) and Pymongo

1. Install [cf-ssh](https://docs.18f.gov/getting-started/cf-ssh/).
2. Use `cf-ssh -f PATH_TO_MANIFEST` command to access application.
3. Once inside you should have access to Python (and Ruby) and can install [Pymongo](https://api.mongodb.org/python/current/installation.html).
  - `git clone git://github.com/mongodb/mongo-python-driver.git pymongo`
  - `cd pymongo`
  - `python setup.py install --user`
4. The client URI used to [connect](https://api.mongodb.org/python/current/tutorial.html) to the database is in the VCAP_SERVICES json. It can be viewed outside of the server using `cf env APP_NAME`

