SERVICE_NAME=reviso-server
CLUSTER_NAME=reviso
deployment_id=$(terraform output -raw ecs_deployment_task_definition)
DEPLOYMENT_STATUS="unknown"
MAX_ATTEMPTS=60
ATTEMPT=0
while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
  DEPLOYMENT_STATUS=$(aws ecs describe-services --cluster ${CLUSTER_NAME} --services ${SERVICE_NAME} --query "services[0].deployments[?taskDefinition=='$deployment_id'].rolloutState" --output text)
  if [ "$DEPLOYMENT_STATUS" == "COMPLETED" ]; then
    echo "Deployment completed successfully."
    break
  elif [ "$DEPLOYMENT_STATUS" == "FAILED" ]; then
    echo "Deployment failed."
    break
  else
    echo "Deployment status of ${deployment_id}: ${DEPLOYMENT_STATUS}. Waiting for completion... Attempt $((ATTEMPT+1))/$MAX_ATTEMPTS"
    sleep 10
  fi
done
