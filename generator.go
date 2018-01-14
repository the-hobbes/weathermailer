package main

import (
  "io/ioutil"
  "log"
  "github.com/golang/protobuf/proto"
  pb "github.com/weathermailer/proto"
)

const FNAME = "proto/conditions.pb"

func MakeThunderstorms() pb.Weather {
  p := pb.Weather{
        Sayings: []*pb.Weather_FolksySaying {
          {Saying: "Blow’n fit to make a rabbit cry",
          Kind: pb.Weather_THUNDERSTORM},
        },
  }
  return p
}

func MakeDrizzles() pb.Weather {
  p := pb.Weather{
        Sayings: []*pb.Weather_FolksySaying {
          {Saying: "Drizzly-drazzly", Kind: pb.Weather_DRIZZLE},
        },
  }
  return p
}

func MakeRains() pb.Weather {
  p := pb.Weather{
        Sayings: []*pb.Weather_FolksySaying {
          {Saying: "Rain placeholder", Kind: pb.Weather_RAIN},
        },
  }
  return p
}

func MakeSnows() pb.Weather {
  p := pb.Weather{
        Sayings: []*pb.Weather_FolksySaying {
          {Saying: "Snow placeholder", Kind: pb.Weather_SNOW},
        },
  }
  return p
}

func MakeAtmospheres() pb.Weather {
  p := pb.Weather{
        Sayings: []*pb.Weather_FolksySaying {
          {Saying: "Capful of wind", Kind: pb.Weather_ATMOSPHERE},
        },
  }
  return p
}

func MakeClears() pb.Weather {
  p := pb.Weather{
        Sayings: []*pb.Weather_FolksySaying {
          {Saying: "Clear placeholder", Kind: pb.Weather_CLEAR},
        },
  }
  return p
}

func MakeClouds() pb.Weather {
  p := pb.Weather{
        Sayings: []*pb.Weather_FolksySaying {
          {Saying: "Clouds placeholder", Kind: pb.Weather_CLOUDS},
        },
  }
  return p
}

func MakeExtremes() pb.Weather {
  p := pb.Weather{
        Sayings: []*pb.Weather_FolksySaying {
          {Saying: "Extremes placeholder", Kind: pb.Weather_EXTREME},
          {Saying: "Thundestorm quote 2 here", Kind: pb.Weather_EXTREME},
        },
  }
  return p
}

func MakeAdditionals() pb.Weather {
  p := pb.Weather{
        Sayings: []*pb.Weather_FolksySaying {
          {Saying: "Additionals placeholder", Kind: pb.Weather_ADDITIONAL},
          {Saying: "Thundestorm quote 2 here", Kind: pb.Weather_ADDITIONAL},
        },
  }
  return p
}

func MakeColds() pb.Weather {
  p := pb.Weather{
        Sayings: []*pb.Weather_FolksySaying {
          {Saying: "Hasn't been this cold since eighteen-hundred-and-froze-" +
                   "to-death", Kind: pb.Weather_COLD},
          {Saying: "It's so cold, milk cows gave icicles",
          Kind: pb.Weather_COLD},
          {Saying: "Colder than a witches tit in a brass bra",
          Kind: pb.Weather_COLD},
        },
  }
  return p
}

func MakeHots() pb.Weather {
  p := pb.Weather{
        Sayings: []*pb.Weather_FolksySaying {
          {Saying: "As snug as a bug in a rug", Kind: pb.Weather_HOT},
        },
  }
  return p
}

func DoGenerateProto() {
  conditions := &pb.WeatherConditions{}
  t := MakeThunderstorms()
  conditions.Weathers = append(conditions.Weathers, &t)
  d := MakeDrizzles()
  conditions.Weathers = append(conditions.Weathers, &d)
  r := MakeRains()
  conditions.Weathers = append(conditions.Weathers, &r)
  s := MakeSnows()
  conditions.Weathers = append(conditions.Weathers, &s)
  a := MakeAtmospheres()
  conditions.Weathers = append(conditions.Weathers, &a)
  c := MakeClears()
  conditions.Weathers = append(conditions.Weathers, &c)
  cl := MakeClouds()
  conditions.Weathers = append(conditions.Weathers, &cl)
  e := MakeExtremes()
  conditions.Weathers = append(conditions.Weathers, &e)
  ad := MakeAdditionals()
  conditions.Weathers = append(conditions.Weathers, &ad)
  co := MakeColds()
  conditions.Weathers = append(conditions.Weathers, &co)
  h := MakeHots()
  conditions.Weathers = append(conditions.Weathers, &h)

  out, err := proto.Marshal(conditions)
  if err != nil {
        log.Fatalln("Failed to encode conditions:", err)
  }
  if err := ioutil.WriteFile(FNAME, out, 0644); err != nil {
        log.Fatalln("Failed to write conditions:", err)
  }
}
