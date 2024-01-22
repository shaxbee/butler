variable "schema_src" {
    type = list(string)
    default = [
        "file://schema",
    ]
}

env "local" {
    src = var.schema_src
    url = "sqlite://local.db"
    dev = "sqlite://dev?mode=memory"
}

env "example" {
    src = var.schema_src
    url = "sqlite://example.db"
    dev = "sqlite://dev?mode=memory"
}