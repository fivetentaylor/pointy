import os
import json
import urllib
import http.client

def lambda_handler(event, context):
    print(json.dumps(event))

    try:
        subject = event["Records"][0]["Sns"]["Subject"]
    except:
        subject = "ECS Alert"


    alert_url = "https://us-west-2.console.aws.amazon.com/ecs/v2/clusters/reviso/services/reviso-server/health?region=us-west-2#"
    message = f"{subject} <{alert_url}|ECS Service Health Dashboard>."
    
    webhook_url = os.environ["SLACK_WEBHOOK_URL"]
    parsed_url = urllib.parse.urlparse(webhook_url)
    
    # Message payload
    payload = json.dumps({'text': message})
    
    # Headers
    headers = {
        'Content-Type': 'application/json',
    }
    
    # HTTP connection
    conn = http.client.HTTPSConnection(parsed_url.netloc)
    conn.request("POST", parsed_url.path, body=payload, headers=headers)
    
    response = conn.getresponse()
    data = response.read()
    conn.close()
    
    return {
        'statusCode': response.status,
        'body': data.decode('utf-8')
    }

if __name__ == "__main__":
    lambda_handler({}, None)

