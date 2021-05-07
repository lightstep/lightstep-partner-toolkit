import logging as log
import json
from botocore.vendored import requests
import os

def main(event, context):
    log.getLogger().setLevel(log.DEBUG)

    log.debug('Input event: %s', event)
    request_type = event['RequestType']

    if request_type == 'Create': return on_create(event, context)
    if request_type == 'Update': return on_update(event, context)
    if request_type == 'Delete': return on_delete(event, context)

    raise RuntimeError('Unknown request type')

def create_dashboard(name, charts, org, project, api_key):
    response = requests.post('https://api.lightstep.com/public/v0.2/%s/projects/%s/metric_dashboards'%(org, project), 
        headers={
            'Authorization': 'Bearer %s'%api_key
        },
        json={
            'data': {
                'attributes': {
                    'name': name,
                    'charts': charts
                }
            }
        })
    return response.json()['data']['id']

def update_dashboard(dashboard_id, name, charts, org, project, api_key):
    payload = {
            'data': {
                'attributes': {
                    'name': name,
                    'charts': charts
                }
            }
        }
    log.debug('update payload:')
    log.debug(json.dumps(payload))
    response = requests.put('https://api.lightstep.com/public/v0.2/%s/projects/%s/metric_dashboards/%s'%(org, project, dashboard_id), 
        headers={
            'Authorization': 'Bearer %s'%api_key
        },
        json=payload)
    return response.json()

def delete_dashboard(dashboard_id, org, project, api_key):
    response = requests.delete('https://api.lightstep.com/public/v0.2/%s/projects/%s/metric_dashboards/%s'%(org, project, dashboard_id), 
        headers={
            'Authorization': 'Bearer %s'%api_key
        })
    return response

def on_create(event, context):
  props = event["ResourceProperties"]
  name = props['name']
  lightstepOrg = props['lightstepOrg']
  lightstepProj = props['lightstepProject']
  charts = json.loads(props['charts'])
  apiKey = os.getenv('LIGHTSTEP_API_KEY')
  log.debug('Creating dashboard "%s" in %s for %s...'%(name, lightstepOrg, lightstepProj))
  physical_id = create_dashboard(name, charts, lightstepOrg, lightstepProj, apiKey)
  attributes = {
    'Response': 'created dashboard "%s"' % physical_id
  }
  return { 'PhysicalResourceId': physical_id, 'Data': attributes }

def on_update(event, context):
  props = event["ResourceProperties"]
  physical_id = event["PhysicalResourceId"]
  lightstepOrg = props['lightstepOrg']
  lightstepProj = props['lightstepProject']
  apiKey = os.getenv('LIGHTSTEP_API_KEY')
  name = props['name']
  charts = json.loads(props['charts'])
  log.debug('Updating dashboard "%s" in %s for %s...'%(physical_id, lightstepOrg, lightstepProj))
  log.debug('charts:')
  log.debug(charts)
  response = update_dashboard(physical_id, name, charts, lightstepOrg, lightstepProj, apiKey)
  attributes = { 'Response': json.dumps(response) }
  return { 'PhysicalResourceId': physical_id, 'Data': attributes }

def on_delete(event, context):
  props = event["ResourceProperties"]
  lightstepOrg = props['lightstepOrg']
  lightstepProj = props['lightstepProject']
  apiKey = os.getenv('LIGHTSTEP_API_KEY')
  physical_id = event["PhysicalResourceId"]
  log.debug('Deleting dashboard "%s" in %s for %s...'%(physical_id, lightstepOrg, lightstepProj))
  delete_dashboard(physical_id, lightstepOrg, lightstepProj, apiKey)
  attributes = {
    'Response': 'deleted dashboard "%s"' % physical_id
  }
  return { 'Data': attributes }
