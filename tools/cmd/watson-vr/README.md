# watson-vr

A comand line tool to manage Watson Visual Recognition.

If you want to use Watson Visual Recognition API in your golang application, see [GoDoc](`https://godoc.org/github.com/speedland/go/services/watson/visualrecognition`)

# Install

# Usage

You need to prepare Watson Visual Recognition API key on [IBM Bluemix](`https://console.ng.bluemix.net/registration/`)
and set `WATSON_API_KEY` environment variable (or pass `-k {key}` option on each command).

## Classify Image

    # From an URL
    $ watson-vr classify https://github.com/watson-developer-cloud/doc-tutorial-downloads/raw/master/visual-recognition/fruitbowl.jpg
    https://github.com/watson-developer-cloud/doc-tutorial-downloads/raw/master/visual-recognition/fruitbowl.jpg:
	    default:
		    - fruit: 0.937027
		    - apple: 0.668188
		    - banana: 0.549834
		    - food: 0.524979
		    - orange: 0.500000

    # From a local file
    $ watson-vr classify ./fruitbowl.jpg
    fruitbowl.jpg:
	    default:
		    - fruit: 0.937027
		    - apple: 0.668188
		    - banana: 0.549834
		    - food: 0.524979
		    - orange: 0.500000

## Detect Faces

    # Show detection result
    $ watson-vr detect-faces https://github.com/watson-developer-cloud/doc-tutorial-downloads/raw/master/visual-recognition/prez.jpg
    https://github.com/watson-developer-cloud/doc-tutorial-downloads/raw/master/visual-recognition/prez.jpg:
	    - location: (256, 64, 92, 159)
          age: 35-35 (0.446989)
          gender: MALE (0.446989)
          identity: Barack Obama /people/politicians/democrats/barack obama (0.970688)

    # Extract face images
    $ watson-vr detect-faces --output ./ https://github.com/watson-developer-cloud/doc-tutorial-downloads/raw/master/visual-recognition/prez.jpg
    $ open ./prez-face-0.jpg

## Create A Cusotom Classifier

You need to prepare the directory that contains example images.

    my_dir/
        /class_name_1/
            a.png
            b.jpg
        /class_name_2/
            c.png
        /.../
        /negative/
            x.png
            y.jpg

Then execute `create-classifier` command to point ./my_dir/

    $ watson-vr create-classifier myclassifier ./my_dir/
    myclassifier_2069560411 created

## List Classifiers

    $ watson-vr list-classifiers
    +-------------------------+--------------+--------+------------------+
    |           ID            |     NAME     | STATUS |     CREATED      |
    +-------------------------+--------------+--------+------------------+
    | myclassifier_2069560411 | myclassifier | ready  | 2016/11/24 11:58 |
    +-------------------------+--------------+--------+------------------+

## Show A Classifier

    $ watson-vr get-classifier myclassifier_2069560411
    myclassifier_2069560411:
    	name: myclassifier
	    status: ready
	    classes:
		    class_name_1
		    class_name_2
            ...

## Update A Classifier

You need to prepare the directory that contains updated images as described `Create A Custom Classifier` section.

    $ watson-vr update-classifier myclassifier_2069560411 ./my_dir/

## Delete A Classifier

    $ watson-vr delete-classifier myclassifier_2069560411
    myclassifier_2069560411 deleted


## Collections

Not supported yet.
