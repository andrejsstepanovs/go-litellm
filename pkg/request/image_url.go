package request

// {
// 	"type": "image_url",
// 	"image_url": {
// 		"url": "https://upload.wikimedia.org/wikipedia/commons/thumb/d/dd/Gfp-wisconsin-madison-the-nature-boardwalk.jpg/2560px-Gfp-wisconsin-madison-the-nature-boardwalk.jpg",
// 		"format": "image/jpeg"
// 	}
// }

type ImageUrl struct {
	URL    string `json:"url"`
	Format string `json:"format"` // Optional field for image format
}
