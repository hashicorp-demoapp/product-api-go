network "local" {
    subnet = "10.5.0.0/16"
}

container "db" {
    network {
        name = "network.local"
    }

    image {
        name = "hashicorpdemoapp/product-api-db:v0.0.11"
    }

    env {
        key = "POSTGRES_DB"
        value = "products"
    }

    env {
        key = "POSTGRES_USER"
        value = "postgres"
    }
    
    env {
        key = "POSTGRES_PASSWORD"
        value = "password"
    }
}

container "api" {
    network {
        name = "network.local"
    }

    image {
        name = "hashicorpdemoapp/product-api:v0.0.11"
    }

    volume {
        source = "./config.json"
        destination = "/config/config.json"
    }

    env {
        key = "CONFIG_FILE"
        value = "/config/config.json"
    }
}

ingress "db" {
    target = "container.db"
    
    network {
        name = "network.local"
    }

    port {
        local = 5432
        remote = 5432
        host = 15432
    }
}

ingress "api" {
    target = "container.api"
    
    network {
        name = "network.local"
    }

    port {
        local = 9090
        remote = 9090
        host = 19090
    }
}