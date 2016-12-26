A little utility that presents current time at given address.

Uses [Google Maps API][maps], so you'll need network access.

[maps]: https://developers.google.com/maps/documentation/timezone

# Example Usage

    $ ./timeat 'los angeles'
    Los Angeles, CA, USA: Fri May 23, 2014 22:33
    $ ./timeat paris
    Paris, France: Sat May 24, 2014 07:37
    Paris, TX, USA: Sat May 24, 2014 00:37
    Paris, TN 38242, USA: Sat May 24, 2014 00:37
    Paris, IL 61944, USA: Sat May 24, 2014 00:37
    Paris, KY 40361, USA: Sat May 24, 2014 01:37
    $

# Installing

You can either `go get github.com/tebeka/timeat` or download the files from
[github][gh].

[gh]: https://github.com/tebeka/timeat/releases
