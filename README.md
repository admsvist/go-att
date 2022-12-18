# go-att

A simple project to use the database of cities in csv format

## API ENDPOINTS

### All Cities

- Path : `/cities`
- Method: `GET`
- Response: `200`

#### Get Query Filters

- Region: `/cities?region=Ростовская область`
- District: `/cities?district=Южный`
- Population Range: `/cities?population=1000000-5000000`
- Foundation Range: `/cities?foundation=1800-2022`

### Create City

- Path : `/cities`
- Method: `POST`
- Fields: `id, name, region, district, population, foundation`
- Response: `201`

### Details a City

- Path : `/cities/{id}`
- Method: `GET`
- Response: `200`

### Update Population City

- Path : `/cities/{id}/population`
- Method: `PATCH`
- Fields: `population`
- Response: `200`

### Delete Post

- Path : `/cities/{id}`
- Method: `DELETE`
- Response: `204`

## Required Packages
- [Go CSV](https://github.com/gocarina/gocsv)
- [chi](https://github.com/go-chi/chi)
