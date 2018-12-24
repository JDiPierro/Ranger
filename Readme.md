# Ranger

Preserving the environment

Ranger loads configuration from the environment. I enjoyed the API of [spf13/viper](https://github.com/spf13/viper) but didn't want to pull in all the extra libraries that Viper uses for parsing files or communicating with external configuration stores. Ranger keeps it simple by only loading configuration from the Environment.

## Installing

The recommended way to install Ranger is with Dep:

```$bash
dep ensure --add github.com/jdipierro/ranger
```

## Usage

Define a struct to hold your configuration. You can use most data types.

```$go
type Config struct {
  HTTPAddr    string
  HTTPPort    string
  TTLMinutes  int
  DBAddress   string
}
```

Give Ranger a default for each key you'd like loaded from the environment, or mark a key as required.

```$go
ranger.SetDefault("HTTPAddr", "0.0.0.0") 
ranger.SetDefault("HTTPPort", "8080")
ranger.SetDefault("TTLMinutes", 60)
ranger.SetRequired("DBAddress")
```

Then unmarshal the environment into your config object:

```
c := new(Config)
err := ranger.Unmarshal(&c)
if err != nil {
  panic(err)
}
```
