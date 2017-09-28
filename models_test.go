package jsonapi

import (
	"context"
	"fmt"
	"time"
)

type key int

const (
	keyBaseURL key = iota
	// ...
)

var testCtx context.Context

func init() {
	testCtx = context.WithValue(context.Background(), keyBaseURL, "http://example.com")
}

type BadModel struct {
	ID int `jsonapi:"primary"`
}

type ModelBadTypes struct {
	ID           string     `jsonapi:"primary,badtypes"`
	StringField  string     `jsonapi:"attr,string_field"`
	FloatField   float64    `jsonapi:"attr,float_field"`
	TimeField    time.Time  `jsonapi:"attr,time_field"`
	TimePtrField *time.Time `jsonapi:"attr,time_ptr_field"`
}

type ModelWithUUIDs struct {
	ID                 UUID                     `jsonapi:"primary,customtypes"`
	UUIDField          UUID                     `jsonapi:"attr,uuid_field"`
	LatestRelatedModel *RelatedModelWithUUIDs   `jsonapi:"relation,latest_relatedmodel"`
	RelatedModels      []*RelatedModelWithUUIDs `jsonapi:"relation,relatedmodels"`
}
type RelatedModelWithUUIDs struct {
	ID        UUID `jsonapi:"primary,relatedtypes"`
	UUIDField UUID `jsonapi:"attr,uuid_field"`
}

type ModelWithUUIDPtrs struct {
	ID                 *UUID                       `jsonapi:"primary,customtypes"`
	UUIDField          *UUID                       `jsonapi:"attr,uuid_field"`
	LatestRelatedModel *RelatedModelWithUUIDPtrs   `jsonapi:"relation,latest_relatedmodel"`
	RelatedModels      []*RelatedModelWithUUIDPtrs `jsonapi:"relation,relatedmodels"`
}

type RelatedModelWithUUIDPtrs struct {
	ID        *UUID `jsonapi:"primary,relatedtypes"`
	UUIDField *UUID `jsonapi:"attr,uuid_field"`
}

type UUID struct {
	string
}

func UUIDFromString(s string) (*UUID, error) {
	return &UUID{s}, nil
}
func (u UUID) String() string {
	return u.string
}

func (u UUID) Equal(other UUID) bool {
	return u.string == other.string
}

type WithPointer struct {
	ID       *uint64  `jsonapi:"primary,with-pointers"`
	Name     *string  `jsonapi:"attr,name"`
	IsActive *bool    `jsonapi:"attr,is-active"`
	IntVal   *int     `jsonapi:"attr,int-val"`
	FloatVal *float32 `jsonapi:"attr,float-val"`
}

type Timestamp struct {
	ID   int        `jsonapi:"primary,timestamps"`
	Time time.Time  `jsonapi:"attr,timestamp,iso8601"`
	Next *time.Time `jsonapi:"attr,next,iso8601"`
}

type Car struct {
	ID    *string `jsonapi:"primary,cars"`
	Make  *string `jsonapi:"attr,make,omitempty"`
	Model *string `jsonapi:"attr,model,omitempty"`
	Year  *uint   `jsonapi:"attr,year,omitempty"`
}

type Post struct {
	Blog
	ID            uint64     `jsonapi:"primary,posts"`
	BlogID        int        `jsonapi:"attr,blog_id"`
	ClientID      string     `jsonapi:"client-id"`
	Title         string     `jsonapi:"attr,title"`
	Body          string     `jsonapi:"attr,body"`
	Comments      []*Comment `jsonapi:"relation,comments"`
	LatestComment *Comment   `jsonapi:"relation,latest_comment"`
}

type Comment struct {
	ID       int    `jsonapi:"primary,comments"`
	ClientID string `jsonapi:"client-id"`
	PostID   int    `jsonapi:"attr,post_id"`
	Body     string `jsonapi:"attr,body"`
}

type Book struct {
	ID          uint64  `jsonapi:"primary,books"`
	Author      string  `jsonapi:"attr,author"`
	ISBN        string  `jsonapi:"attr,isbn"`
	Title       string  `jsonapi:"attr,title,omitempty"`
	Description *string `jsonapi:"attr,description"`
	Pages       *uint   `jsonapi:"attr,pages,omitempty"`
	PublishedAt time.Time
	Tags        []string `jsonapi:"attr,tags"`
}

type Blog struct {
	ID            int       `jsonapi:"primary,blogs"`
	ClientID      string    `jsonapi:"client-id"`
	Title         string    `jsonapi:"attr,title"`
	Posts         []*Post   `jsonapi:"relation,posts"`
	CurrentPost   *Post     `jsonapi:"relation,current_post"`
	CurrentPostID int       `jsonapi:"attr,current_post_id"`
	CreatedAt     time.Time `jsonapi:"attr,created_at"`
	ViewCount     int       `jsonapi:"attr,view_count"`
}

func (b *Blog) JSONAPILinks(ctx context.Context) *Links {
	return &Links{
		"self": fmt.Sprintf("https://example.com/api/blogs/%d", b.ID),
		"comments": Link{
			Href: fmt.Sprintf("https://example.com/api/blogs/%d/comments", b.ID),
			Meta: Meta{
				"counts": map[string]uint{
					"likes":    4,
					"comments": 20,
				},
			},
		},
	}
}

func (b *Blog) JSONAPIRelationshipLinks(ctx context.Context, relation string) *Links {
	if relation == "posts" {
		return &Links{
			"related": Link{
				Href: fmt.Sprintf("https://example.com/api/blogs/%d/posts", b.ID),
				Meta: Meta{
					"count": len(b.Posts),
				},
			},
		}
	}
	if relation == "current_post" {
		return &Links{
			"self": fmt.Sprintf("https://example.com/api/posts/%s", "3"),
			"related": Link{
				Href: fmt.Sprintf("https://example.com/api/blogs/%d/current_post", b.ID),
			},
		}
	}
	return nil
}

func (b *Blog) JSONAPIMeta(ctx context.Context) *Meta {
	return &Meta{
		"detail": "extra details regarding the blog",
	}
}

func (b *Blog) JSONAPIRelationshipMeta(ctx context.Context, relation string) *Meta {
	if relation == "posts" {
		return &Meta{
			"this": map[string]interface{}{
				"can": map[string]interface{}{
					"go": []interface{}{
						"as",
						"deep",
						map[string]interface{}{
							"as": "required",
						},
					},
				},
			},
		}
	}
	if relation == "current_post" {
		return &Meta{
			"detail": "extra current_post detail",
		}
	}
	return nil
}

type BadComment struct {
	ID   uint64 `jsonapi:"primary,bad-comment"`
	Body string `jsonapi:"attr,body"`
}

func (bc *BadComment) JSONAPILinks(ctx context.Context) *Links {
	return &Links{
		"self": []string{"invalid", "should error"},
	}
}
