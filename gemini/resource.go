package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/api/iam/v1"
)

type configer interface {
	ReadStruct(string, any) error
}

type config struct {
	ApiKey      string
	ApiEndpoint string
	ProjectID   string
	LocationID  string
	ModelID     string
}

// func CallAPI() string {

// 	model := client.GenerativeModel("gemini-pro-vision")
// 	img := genai.ImageData("jpeg", image_bytes)
// 	prompt := genai.Text("Please give me a recipe for this:")
// 	resp, err := model.GenerateContent(ctx, img, prompt)
// }

func GenerateContentFromText(w io.Writer, cfgr configer) error {
	ctx := context.Background()

	s, err := iam.NewService(ctx)
	if err != nil {
		return fmt.Errorf("iam.NewService: %w", err)
	}
	_ = s
	cfg := config{}
	if err := cfgr.ReadStruct("gemini", &cfg); err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}

	client, err := genai.NewClient(ctx, cfg.ProjectID, cfg.LocationID)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	gemini := client.GenerativeModel(cfg.ModelID)
	prompt := genai.Text(
		`Topic:	why my patients do not follow my diet change recommendations survey
		 Number of questions: 22
		 Areas of survey: 
			1. perceived importance of the change
        	2. self efficacy problems
			3. bad eating habits under stress
			4. bad eating due to succumbing to environmental influences
			5. eating as a compensation to dietary restrictions
			6. other psychological problems leading to eating disorders
		Question distribution:
			for areas 1-5: 4 questions per area
			for area 6: 2 questions
			Language: Polish`)

	gemini.ResponseMIMEType = `application/json`
	gemini.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(`jestes specjalista od psychologii jedzenia.`),
			genai.Text(`twoja dziedzina to zmiana nawykow, w szczegolnosci nawykow zwiazanych z odzywianiem.`),

			genai.Text(`twoj odbiorca to specjalista, glownie dietetyk pracujacy z pacjentem.`),
			genai.Text(`unikaj dzielenia jedzenia na zdrowe i niezdrowe, zamiast pojecia jedzenia niezdrowego uzywaj pojecia  jedzenia o wysokiej gestosci energetycznej oraz duzej smakowitosci (slodycze, przekaski itp).`),
			genai.Text(`wez pod uwage ze niewlasciwe zachowania zywieniowe moga wynikac z przyczyn fizjologicznych (np: restrykcje dietetyczne) lub psychologicznych.`),
		},
	}

	resp, err := gemini.GenerateContent(ctx, prompt)
	if err != nil {
		return fmt.Errorf("error generating content: %w", err)
	}
	// See the JSON response in
	// https://pkg.go.dev/cloud.google.com/go/vertexai/genai#GenerateContentResponse.
	rb, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}
	fmt.Fprintln(w, string(rb))
	return nil
}
