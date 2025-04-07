package dbstore

import (
	"sjb_site/internal/store"

	"gorm.io/gorm"
)

type GroupStore struct {
	db *gorm.DB
}

type NewGroupStoreParams struct {
	DB *gorm.DB
}

func NewGroupStore(params NewGroupStoreParams) *GroupStore {
	return &GroupStore{
		db: params.DB,
	}
}

func (s *GroupStore) CreateGroup(group *store.Group) error {
	return s.db.Create(group).Error
}

func (s *GroupStore) GetGroup(groupId string) (*store.Group, error) {

	var group store.Group
	err := s.db.Where("id = ?", groupId).First(&group).Error

	if err != nil {
		return nil, err
	}
	return &group, err
}

func (s *GroupStore) GetGroupsByType(groupType string) (*[]store.Group, error) {
	var groups []store.Group
	err := s.db.Where("group_type = ?", groupType).Find(&groups).Error

	if err != nil {
		return nil, err
	}
	return &groups, err
}

func (s *GroupStore) GetJaarclubs(jaarlaag int) (*[]store.Group, error) {
    var groups []store.Group
    subQuery := s.db.Table("groups").Select("YEAR(start_date) as year").Where("group_type = ?", "jaarclub").Group("year").Order("start_date DESC").Offset(jaarlaag).Limit(1)
	err := s.db.Where("group_type = ? AND YEAR(start_date) = (?)", "jaarclub", subQuery).Find(&groups).Error
    return &groups, err
}

func (s *GroupStore) GetCommissies() (*[]store.Group, error) {
	var commissies []store.Group

	query := `
		SELECT 
			g.id, 
			IF(p.name IS NOT NULL, p.name, g.name) AS name,
			g.image,
			g.group_type,
			g.start_date,
			g.end_date
		FROM groups g
		LEFT JOIN parent_groups pg ON pg.child_id = g.id
		LEFT JOIN groups p ON pg.parent_id = p.id
		INNER JOIN (
			SELECT 
				COALESCE(pg.parent_id, g.id) AS group_key,
				MAX(g.start_date) AS max_start
			FROM groups g
			LEFT JOIN parent_groups pg ON pg.child_id = g.id
			WHERE g.group_type = ?
			GROUP BY group_key
		) latest ON latest.max_start = g.start_date 
		         AND COALESCE(pg.parent_id, g.id) = latest.group_key
		WHERE g.group_type = ?
		ORDER BY g.end_date ASC, name;
	`

	err := s.db.Raw(query, "commissie", "commissie").Scan(&commissies).Error
	if err != nil {
		return nil, err
	}
	return &commissies, nil
}

func (s *GroupStore) GetSimilarGroups(group *store.Group) (*[]store.Group, string, error) {
    var groups []store.Group
    var title string
    var err error
    switch group.GroupType {   
    case "jaarclub":
        err = s.db.Where("group_type = ? AND year(start_date) = year(?) AND id != ?", "jaarclub", group.StartDate, group.ID).Find(&groups).Error
        title = "Jaarclubs uit " + group.StartDate.Format("2006")
    case "barploeg":
        err = s.db.Where("group_type = ? AND end_date IS NULL AND id != ?", "barploeg", group.ID).Find(&groups).Error
        title = "Barploegen"
    case "commissie":
        subquery := s.db.Select("parent_id").Table("parent_groups p").Where("p.child_id = ?", group.ID)
        err = s.db.Joins("JOIN parent_groups pg ON groups.id = pg.child_id").
                Where("pg.parent_id = (?) AND groups.id != ?", subquery, group.ID).
                Find(&groups).Error
        title = "Commissies"
    case "bestuur":
        err = s.db.Where("group_type = ? AND id != ?", "bestuur", group.ID).Find(&groups).Error
        title = "Besturen"
    case "huis":
        err = s.db.Where("group_type = ? AND id != ?", "huis", group.ID).Find(&groups).Error
        title = "Huizen"
    case "gilde":
        err = s.db.Where("group_type = ? AND end_date IS NULL AND id != ?", "gilde", group.ID).Find(&groups).Error
        title = "Gilden"
    case "werkgroep":
        err = s.db.Where("group_type = ? AND end_date IS NULL AND id != ?", "werkgroep", group.ID).Find(&groups).Error
        title = "Werkgroepen"
    }

    if err != nil {
        return nil, "", err
    }

    return &groups, title, err
}







