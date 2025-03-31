package models

import "time"

type SkillType string

const (
	Language   SkillType = "Language"
	Frameworks SkillType = "FrameWorks"
	Tools      SkillType = "Tools"
)

// 技術などの基本構造遺体
type Skills struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      SkillType `json:"type"`
	HasImage  bool      `json:"has_image"`
	CreatedAt time.Time `json:"created_at"`
}
