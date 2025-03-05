# code

## Introduction


## Key Packages

- `github.com/99designs/gqlgen` : For generating GraphQL servers in Go.
- `github.com/charmbracelet/log` : A charming logger for all your logging needs.
- `github.com/go-chi/chi/v5` : Lightweight and feature-rich router for building Go HTTP services.
- `github.com/redis/go-redis/v9` : A Redis client for Golang.
- `gorm.io` : A developer-friendly ORM for handling interactions with your PostgreSQL database.

## Getting Started

### Prerequisites

- Go 1.18 or higher
- Docker

### Installing

1. Fork the repository, and clone it to your machine

```sh
git clone https://github.com/teamreviso/code
```

2. Move into the project directory.

```sh
cd code
```

3. Download the required Go dependencies.

```sh
make install
```

4. Setup local certs for tls
```sh
make certs
```

5. Setup your database and fill the required information in the `.env` file. Look at the `.env.example`.

6. Run the server locally.

```sh
make dev
```

Now, your server should be running at `localhost:8080`. (or what ever ADDR you set in your `.env` file.

## Deployment

You can build the project using the standard Go build tool. This will create a binary file that can be executed.

```sh
go build -o main .
```

## License

This project is licensed under the MIT License - see the `LICENSE.md` file for details.

## Acknowledgments

This project wouldn't be possible without these wonderful projects and their contributors:

- [GQLGen](https://github.com/99designs/gqlgen)
- [Charm Log](https://github.com/charmbracelet/log)
- [Chi](https://github.com/go-chi/chi)
- [Go-Redis](https://github.com/redis/go-redis)
- [GORM](https://gorm.io)

Please feel free to contribute to this project, report bugs and issues, and suggest improvements.


## Deployment

### Manually deploy a preview branch of the app server:

Login to ECR:

```sh
assume-role staging:admin aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 533267310428.dkr.ecr.us-west-2.amazonaws.com
```

Set environment variables:

```sh
export PR_NUMBER=(your PR number)
export NODE_ENV=production
export WEB_HOST=https://pr-$PR_NUMBER-www.reviso.biz
export APP_HOST=https://pr-$PR_NUMBER-app.reviso.biz
export WS_HOST=wss://pr-$PR_NUMBER-app.reviso.biz
```

Build and push the server image to ECR:

```sh
docker build -t 533267310428.dkr.ecr.us-west-2.amazonaws.com/reviso-server:$(git rev-parse HEAD) -f cmd/reviso/Dockerfile --build-arg NODE_ENV=$NODE_ENV --build-arg APP_HOST=$APP_HOST --build-arg WS_HOST=$WS_HOST --build-arg WEB_HOST=$WEB_HOST .
docker push 533267310428.dkr.ecr.us-west-2.amazonaws.com/reviso-server:$(git rev-parse HEAD)
```

Run the terraform:

```sh
cd terraform/postbuild/preview
```

Set the environment variables:

```sh
export TF_VAR_web_sha=$(git rev-parse HEAD)
export TF_VAR_server_sha=$(git rev-parse HEAD)
export TF_VAR_pr_number=$PR_NUMBER
export TF_VAR_slack_webhook_url=(get from https://api.slack.com/apps/A06KB3LHGAY/incoming-webhooks)
```

```sh
assume-role staging:admin terraform init
assume-role staging:admin terraform plan
assume-role staging:admin terraform apply
```

### Manually deploy a staging branch of the app server:

Login to ECR:

```sh
assume-role staging:admin aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 533267310428.dkr.ecr.us-west-2.amazonaws.com
```

Set environment variables:

```sh
export NODE_ENV=production
export WEB_HOST=https://www.reviso.biz
export APP_HOST=https://app.reviso.biz
export WS_HOST=wss://app.reviso.biz
```

Build and push the server image to ECR:

```sh
docker build -t 533267310428.dkr.ecr.us-west-2.amazonaws.com/reviso-server:$(git rev-parse HEAD) -f cmd/reviso/Dockerfile --build-arg NODE_ENV=$NODE_ENV --build-arg APP_HOST=$APP_HOST --build-arg WS_HOST=$WS_HOST --build-arg WEB_HOST=$WEB_HOST .
docker push 533267310428.dkr.ecr.us-west-2.amazonaws.com/reviso-server:$(git rev-parse HEAD)
```

Run the terraform:

```sh
cd terraform/postbuild/staging
```

Set the environment variables:

```sh
export TF_VAR_web_sha=$(git rev-parse HEAD)
export TF_VAR_server_sha=$(git rev-parse HEAD)
export TF_VAR_slack_webhook_url=https://hooks.slack.com/services/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
```

```sh
assume-role staging:terraform terraform init
assume-role staging:terraform terraform plan
assume-role staging:terraform terraform apply
```

View the logs:

```sh
assume-role staging:terraform aws logs tail --follow /ecs/reviso-server/main
```


### Manually push local images to ECR

#### Build and push the server image to Production ECR

```
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 998899136269.dkr.ecr.us-west-2.amazonaws.com
docker build -t 998899136269.dkr.ecr.us-west-2.amazonaws.com/reviso-server:$(git rev-parse HEAD) -f cmd/reviso/Dockerfile .
```

#### Build and push the web image to Production ECR

> Note you need to remove your node_modules first, since it needs to build from scratch 

```
rm -rf frontend/web/node_modules
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 998899136269.dkr.ecr.us-west-2.amazonaws.com
docker build -t 998899136269.dkr.ecr.us-west-2.amazonaws.com/reviso-server:$(git rev-parse HEAD) -f cmd/reviso/Dockerfile .
```

```sh
aws ecs execute-command  \
    --region <region> \
    --cluster <cluster-name> \
    --task <task-id> \
    --container <container-name> \
    --command "/bin/sh" \
    --interactive
```

#### Pushing up new faktory version


```
docker pull contribsys/faktory:latest
assume-role staging:admin aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 533267310428.dkr.ecr.us-west-2.amazonaws.com
docker tag contribsys/faktory:latest 533267310428.dkr.ecr.us-west-2.amazonaws.com/faktory-cache:latest
docker push 533267310428.dkr.ecr.us-west-2.amazonaws.com/faktory-cache:latest
```

> Note: using staging here


