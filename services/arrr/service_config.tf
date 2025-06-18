locals{
  services = {
    1 = {
      service_name = "prowlarr"
      image_name = "linuxserver/prowlarr:1.37.0"
      username = "default_username"
      password = "ProwlarrPassword1!"
      domain_name = "prowlarr.mallett.family"
      service_port = 9696
      ip_address = "192.168.5.13"
      env = [
        "PUID=1000",
        "PGID=1000",
        "TZ=America/Chicago",
      ]
    },
    2 = {
      service_name = "radarr"
      image_name = "linuxserver/radarr:5.26.2"
      username = "default_username"
      password = "RadarrPassword1!"
      domain_name = "radarr.mallett.family"
      service_port = 7878
      ip_address = "192.168.5.14"
      env = [
        "PUID=1000",
        "PGID=1000",
        "TZ=America/Chicago",
      ]
    },
    3 = {
      service_name = "sonarr"
      image_name = "linuxserver/sonarr:4.0.14.2939-ls281"
      username = "default_username"
      password = "SonarrPassword1!"
      domain_name = "sonarr.mallett.family"
      service_port = 8989
      ip_address = "192.168.5.15"
      env = [
        "PUID=1000",
        "PGID=1000",
        "TZ=America/Chicago",
      ]
    },
    4 = {
      service_name = "lidarr"
      image_name = "linuxserver/lidarr:2.12.4"
      username = "default_username"
      password = "LidarrPassword1!"
      domain_name = "lidarr.mallett.family"
      service_port = 8686
      ip_address = "192.168.5.16"
      env = [
        "PUID=1000",
        "PGID=1000",
        "TZ=America/Chicago",
      ]
    },
    5 = {
      service_name = "readarr"
      image_name = "linuxserver/readarr:0.4.17-develop"
      username = "default_username"
      password = "ReadarrPassword1!"
      domain_name = "readarr.mallett.family"
      service_port = 8787
      ip_address = "192.168.5.17"
      env = [
        "PUID=1000",
        "PGID=1000",
        "TZ=America/Chicago",
      ]
    }
  }
}