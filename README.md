# sitebuild
Sitebuild static site generator


## Example sitebuild.json configuration

```json
{
  "sourceDir": "src",
  "pages": {
    "index.html": {
      "sources": ["layout.html", "home.html"],
      "url": "/"
    },
    "about.html": {
      "sources": ["layout.html", "about.html"],
      "url": "/about-us"
    }
  }
}
```

The pages object is a map of output pages. Each key is the output file name, and the value is an object with the following properties:

- sources: an array of source files to be concatenated together to create the output file

- url: the URL path for the output file

The basepath property specifies the directory where the source files are located. The source files are expected to be in the basepath directory, and the output files will be written a directory named `dist`.
