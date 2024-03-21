package handlers

// Represent shorten endpoint request's input data
type URLShortenInput struct {
	URL string `json:"url"`
	TTL string `json:"ttl"`
}

// App configurations
type Config struct {
	Port      string `default:"8080"`
	Region    string `default:"eu-west-1"`
	RedisHost string `default:"localhost:6379"`
}
