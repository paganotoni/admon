## Admon

Admon is an administration framework for Buffalo application. It aims to provide a way to generate admin portals for buffalo in an easy way. 

![Admon Preview](https://user-images.githubusercontent.com/645522/66008895-eaa15f80-e47d-11e9-840d-11cd0718f2ef.png)


Amon is heavily inspired by [Active admin]("https://activeadmin.info").

### Using Admon

Admon assumes go modules. To install do:

```
go get -u -v github.com/paganotoni/admon
```

Once you've installed the package In your buffalo `app.go` file you can add:

```go
admon.Register(models.MyModel{})
admon.MountTo(app)
```

And assuming `models.MyModel` is a generated buffalo model, it should just work. You could visit `https://localhost:3000/admin/`, and see the Admon dashboard.

### Advanced configuration

The following example shows a bit more of how to customize Admon list table and forms.

```go
admon.Register(models.Team{}).WithOptions(admon.Options{
    Fields: admon.FieldOptionsSet{
        {Name: "Name"}, //Selecting specific columns and order.
        {Name: "ShortName"},
        {Name: "Nickname"},
        {
            Name:     "Gender",
            // Specifying how to render values for this field.
            Renderer: func(value interface{}) *tags.Tag { 
                badgeClass := "badge badge-danger"
                if fmt.Sprintf("%v", value) == "MENS" {
                    badgeClass = "badge badge-primary"
                }

                return tags.New("span", tags.Options{"class": badgeClass, "body": value})
            },
            // Specifying what kind of field will be shown in the form.
            Input: admon.InputOptions{
                Type:          admon.InputTypeSelect,
                SelectOptions: []string{"", "MENS", "WOMENS"},
            },
        },
    },
})
```

### ⚠️ Important

Admon is still under heavily development and still not production ready.





