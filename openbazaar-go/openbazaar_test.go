package openbazaar

import (
	"net/http"
	"testing"
)

const (
	img = "iVBORw0KGgoAAAANSUhEUgAAAfQAAAH0CAIAAABEtEjdAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAADtVJREFUeNrs3T12ItcWBtCW7cAhypypGIFKI1BpBngEjUZgNAKjEVR36Aj1CEAjAGXOQJkzUOYMlNkR77bw0uvVP2oK/d26tXegpX6r7Qenyh+HU7duvXkDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAERqTwloiKIotvlrq9VqNpspF8IdXl/2ifDH4+Pj8LPVauV5/sh/cwj6EPebX25vbzfR7wMA4Q5PL0R2CPHw8/DwcPPLq7yMTcQvFoubm5vJZLK44+gg3GFboQcviiKEeGjJt5yuvJaQ8iHxr6+vZ3ccO4Q7fCXQN2n+Wr35k2T91dXV5I5jinCnubIs63Q6IdPDz8TemqAHmtin93q96XS6boDlcjkcDrvd7ubyL0CCQocekm7dVOHzrN/v13fuBPB5qx5CbT6fr7kTSlGWZXrDKKApsiwbDAbL5VKgf2toE+oj5YGaxbr4rpTyJjZA1EOYsizl9WMmNq6+AnHp9/uGME919bXb7YZPSicV8JqKonDJ9DkMBoPI79QFzGHYfVwTvhUZ1wAa9jQNh0ONPPC8Qi8pbV+rke/1eibywNOPYhqyf0D8CyjNarBxGE8jz/PxeKxtjMRisWi32+rQZD8oAY/X7XZDzy7Z43F+fq4IwKMYskc4mfFBy09KwGMMBoPQtqtDVEaj0ea5rwA7Jrs2OUK2owEke2rG47GTkzcuqGIak5gPHz4oAm8shWQHZVn2ej11iNBqtdrf31cHdO5UFhp2yR6t9+/fKwI6dyorisJIN2btdnuxWKgDOncqyLJsOByqQ7RGo5FkR7hTWUh2t8bEzEwG4U5lZVlaPR2z0LNPJhN1QLhTQVEULqJq26kXF1T5js1GvraQjdlqtWq327YcQOdOBb///rtkj5zNZNC5U02e56FtV4fIHR0dzWYzdUDnzrbKslSEyM3uqAPCnW11Oh2PXY6fS6l8lbEM3zSfz03bI2czGXTuVNPtdiV7/C4uLhQBnTva9tTYTAadOxV0Oh3JHj+bySDcqea3335ThPh5LgcPMJbhc6Fnn8/n6hC50LO32211QOeOtj0pVkCic6ea5XKZ0ta+s9lstVqFn7e3t/d//Ozv5Hl+/5YPDw83v0e+xn9/f9+WAzzgJyXgU51Op+7JHiJvMplcXV2FHN9yF9wH/lqI+FCQkP4HBwdZlkWS+BcXF5IdqGAwGKxrK7z48OH03CUKQd/tdsuynE6nr/VO3TkMVJ7J1C7Tw2vu9/uv9YUj5Gz4fx+Pxy/2fm3lBlSeyYj1RxYwdPTz+fxZ33X43uBcBSoIwVSjZA/NcrR3WoUX1uv1nqOdT+xyN/ASXnGIXFVd9iIOQRwa7eFw2LQ3DkQUQ3VJ9jrOJZ4q5W0LAVRTl4F73SfOIeV7vd5uX5LG47ETFaim3+9L9pcUevBQ80pXX19goSeQmpdczLfzMvYkK18URXhr312EasMfYBeRr3BPfnH3Zij/wEdsaPOdpUDlEUHkbXue5805FmVZfvlZawUksMtkwEAmNp828s2sAPBYkV9NbfL6v00jbzMZqrIrJLG7uLho8sPkwns/OztzGlCVh3Xw0fHxcbSv7fLy0gEC4U5SVqvVaDRSBxDu7CLatSiSHYQ7u4t2md319bWjA8Kd1MxmM0UA4U5SbbtwB+HO7mK++dNjoEG4AyDcid5kMlEEEO4ACHeiZ0MVEO4ACHdqwibmINxJUHOe0QHCHeEOCHe+J+bd0mPeixhitqcEBOv1OtrXtr+/7z5V0LmTmk6nowgg3NlFzK3x27dvHSAQ7uwi5s0XizuOEQh3UlOWpSJAJT8qAUG73Y65O/7ll19ub2///PNPRwp07lQQ/3KU0Lxb8w7CnWpq8cCj4XBoNwKACkJorutgOp3Kd4AK5vO5fIdkGMvwn7o8ijrP8/F4LN9BuLOV6+vrurxU+Q6wraIo1rWyXC6tnwH4vhCX67rp9XoOHMBDhsPhuoYskQR4SLfbXddT+M5h80iAr8uybF1n4/E4vAXHEeBz0+l0XXNlWZrSACQymflsSuNCK8D/1WUfgm3M5/PwWeWYAnw0GAzWCRmPx571AfDx/s91ckQ8wJsQhesUiXig0Wq3FUElg8HAiklA8y7iATTvIh5A8y7iATTvz7sBmcutQCMktubdihqAj1qtVh03eRfxAN/R6/XWTSXigZQ158rqV02nU3vUAAnKsqyZwxnbkAGGMyIeoIZq+oRVEQ/wkFarFUJNsot4IDVJ7gYs4gESeQ7fc6yosWgSqLcG3rZqXTzQCA1f+W4bMiBNrVZrOp3K8W9ZLpf9fj9UyakCyHfXWgEikOe5O1e3GcSHQjlbAPmeoLIsTWkA+Z7mlKbT6ThhgDrlu/n7lobDoRYeqA3XVyutpdHCA3XKd+vftfBAmty/qoUH0mT/GQtpgDTleW5/4Er7jlkLD9RD6EY936PSiMbtrEBt9Ho9q+ArbTrmnAFqM6KxSrLSiMYIHqiNfr8vuLcf0RjBA3Vq4S2Et0oSSJMp/PZcYgXqpNVquddJvgNpKorChVZLaIA0hc7UlEa+8+R+VAJe12w2++OPP/7999/QyKvGt+R5vre3N5lMlAKomSzLDOLN34FkI95ySfkOpKkoChHv/iZAxDcr3+1PAIj4BIWCODEAEZ+gfr/vrADSiXgrau4ZvgNJsWjyfnNgJwMg4g1nAOoT8SHgmryBQaiA0wBIU6vVamzEWzkDiPg02ZYHEPGurAKIeHvOAIj41zKfzx1uQMSbvAOIeMtmAGKW9q1P1rwDIj7BiC/L0sEFmi69nSaXy6XDCvBRp9OZz+fJ5Ht4O44pwY9KQMP99ddf79+/39vby/P8559/rvvb+eeffy4vLx1WgP9kWTYcDk1mANKc0tR9uaQF7wBf0Wq1at3CWzMD8E3dbremLbx9xAAekud5CMo65nv48uHwATw0oqnjWngLIvlBCeABq9Xq5OTk4uKidt85HDuA76vXdgU2EQNIMN+tdgdIM99dUwWooC5L4N3K1HAuqEI1p6ens9ks/tdpb3fhDlSwWq1Cvgt3ImdXSKjs77//3tvbi3zusVgsbA8JUFnku8BbDdlwxjKwo/Pzc0UgWntKAI9p3qMdba9Wq/39fcdI5w5U9uHDh2hfm3Xuwh3Y0Wg0UgTiZCwDjxLzZGZvz3/gOndgJ5PJRBEQ7pCam5sbRUC4g84dhDsAwh0A4Q4g3AEQ7gAId6grd/kj3CFBeZ7H+cJq8bgohDtE6uDgIM4XtlqtHB3hTkMtl8t+v2+wkGTnDjTX/VN7BoOBR27uIHwuRvskpvCx7QDp3Gm6brc7n8/H43Gn01GN7SkXEHXn/qmQ8r1ez6xmG8PhMNrOPfLndwMvHe6fzmoMlB+QZVnMD8h27EC4P2Q6nXa7XY38l/r9fszh7gCBcP++5XKpkf9U+LQLNYk22cNHsmMEwr1aamjkg/BRF3PbHl6e0xuEe2WbRr6xl+w6nc46br1ez+kNwn138/m83+83ao18eLMxD2RcTQXePO2QtwnjmvAGwzuNPNnDZ49zG4T7ExsOh6mmfC2SfXMInNsg3J8xYlJK+TzPa5HsQSi7cxuE+0s0kiFuaj2XD68//jn7PdsEgXB/6cXX/X6/Xtf6QlDGvMeAFe7A64f7p1f84m/nW61W5PegWgQJxBXuny2mHAwGIZLi6ejDR054STWaw5jJ8KU9JWh4uEf1elar1Ww2u7q6mt1ZLBYvnOmdTuft27f1XSQeinZ0dOTERrgL96i3l9pkfXBzc7PJ+ieP+xDoIcqPj4+Lokjgxp/T09OLiwsnNsJduNdv78CQ8pvQv7293fyy+d8nk8kD/1Sr1dpkd3bn4OBgE+sprccP1Wi32x6dysZPSkC9bDLakyi+FHp2yY7Onbp27nxLaNtf+CoFMfMMVUikbZfs6NzRuWvb0bkD2nZ07ujc0bajcwce5fz8XLKjc0fnnpQQ60dHR1ZAonOHpJydnUl2hDskZXRHHfgqY5lGM5apL5sNoHOHBJ2enkp2hDsk5d27dwYyPMxYptGMZerIpu0Id4R7aoza2ZKxDNQp2U9OTiQ7wh2Scnp6ev9wEhDukEiyu4iKcIfUkt3DURHuINkR7oBkR7gDL2O1Wkl2dvaTEkCcyX5ycmJtDDp3SMfmHlTJjnCHdIxGo9Cze7gSwh3ScXZ29uuvv7oHlcczc4cohFY9xLpRDDp3SMe7d+8M2YEnMx6P17yq+XxeFIVTEXhiIVlE/KtYLpf9ft8ZCDyjPM8Hg4HAfTHD4TDLMice8BJC3IRecj6fC9/nE74nmcMAr6Pb7ZrViHVAI49YB+pmM5FfLpdiuqpQt1A9pxAQtU6nI+W3XODY6/VarZZzBqhZL1+WpYnNl6sbtepACrIs63a7w+Gwye38JtPD1xrnA/HYUwKesJ0viuL4+Dj8bMJEYjabXV5ejkYj2wYg3GlQ0Ach6De/pBTok8nk6uoq/LR3I8Kdpgu9fJZlh4eHm6yvUV+/WCxCoF9fX0/uOJQId/imEO4h4rM7IfHDHyNZBh5yPPTjoTFf3JHmCHd4ApumfhP64Y+b3L//MHiqTnzzkKMQ4qEfvw/0zU+HAOEOr9z4b//3NeAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMTjfwIMAPllFq5VUEcDAAAAAElFTkSuQmCC"
)

func TestLogin(t *testing.T) {
	_, err := setupClient()
	if err != nil {
		t.Fatalf("Failed to login.")
	}
}

func TestUploadImage(t *testing.T) {
	client, _ := setupClient()

	hash, err := UploadImage(client, img)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(hash) == 0 {
		t.Fatalf("Failed to upload image. Hash is zero length.")
	}
}

func TestContract(t *testing.T) {
	client, _ := setupClient()
	hash, _ := UploadImage(client, img)
	contract := Contract{
		Title:          "Test Title",
		Image:          hash,
		Price:          "34.50",
		Tags:           []string{"test1", "test2"},
		Description:    "Test Description",
		ShippingOrigin: "ALL",
		CurrencyCode:   "USD",
	}
	err := PostContract(client, contract)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

// Setup the client and login.
func setupClient() (client *http.Client, err error) {
	username := "username"
	password := "password"
	return Client(username, password)
}
