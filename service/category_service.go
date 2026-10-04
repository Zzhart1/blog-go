package service

import (
	"errors"

	"blog/model"
)

func CreateCategory(name string) (*model.Category, error) {
	cate := model.Category{Name: name}
	if err := model.DB.Create(&cate).Error; err != nil {
		return nil, err
	}
	return &cate, nil
}

func GetCategoryList() ([]model.Category, error) {
	var list []model.Category
	err := model.DB.Order("id ASC").Find(&list).Error
	return list, err
}

// UpdateCategory 更新分类名
func UpdateCategory(idStr, name string) (*model.Category, error) {
	id, err := parseID(idStr)
	if err != nil {
		return nil, err
	}
	var cate model.Category
	if err := model.DB.First(&cate, id).Error; err != nil {
		return nil, err
	}
	if err := model.DB.Model(&cate).Select("Name").Updates(model.Category{Name: name}).Error; err != nil {
		return nil, err
	}
	return &cate, nil
}

// DeleteCategory 删除分类；分类下仍有文章时拒绝删除
func DeleteCategory(idStr string) error {
	id, err := parseID(idStr)
	if err != nil {
		return err
	}
	var cate model.Category
	if err := model.DB.First(&cate, id).Error; err != nil {
		return err
	}
	var count int64
	if err := model.DB.Model(&model.Article{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该分类下仍有文章，无法删除")
	}
	return model.DB.Delete(&cate).Error
}
