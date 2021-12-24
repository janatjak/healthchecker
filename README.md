# Health checker

Simple health checker (status source) for atlassian status page. 

## Install

1. Visit https://www.atlassian.com/software/statuspage, create new page and add components.
2. Create API key in https://manage.statuspage.io/organizations/<page-id>/api-info
3. Create docker-stack.yml file

```yaml
version: '3.7'
services:
    checker:
        image: janatjak/healthchecker
        environment:
            CONFIG: >
                {
                    "apiKey": "<status page key>",
                    "pageId": "<status page id>",
                    "mainComponentId": "<main component id>",
                    "components": [
                        {
                            "id": "<component id>",
                            "url": "<url>"
                        }
                    ]
                }
```

4. Create docker stack
```shell
docker stack deploy -c docker-stack.yml status
```