package gormutils

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
)
var defaultTimezone = "America/Chicago"
// ConvertDateToUnixTime convert updated_at, created_at, deleted_at to unix timestamp
func ConvertDateToUnixTime(timezone ...string) {
	var tz = defaultTimezone
	if len(timezone) > 0 {
		tz = timezone[0]
	}

	gorm.NowFunc = func() time.Time {
		var now = time.Now()
		loc,err:= time.LoadLocation(tz)
		if err == nil {
			now = now.In(loc)
		}

		return now
	}
	
	DBInstance.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	DBInstance.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	DBInstance.Callback().Delete().Replace("gorm:delete", deleteCallback)
	gorm.DefaultCallback.Create().Before("gorm:save_before_associations").Register("app:update_xid_when_create", updateIDForCreateCallback)\

}

func updateIDForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		if id, ok := scope.Get("ID"); ok {
			if id, ok := id.(string); ok {
				if id == "" {
					scope.SetColumn("ID", xid.New().String())
				}
			}
		} else {
			scope.SetColumn("ID", xid.New().String())
		}

	}
}

// updateTimeStampForCreateCallback will set `CreatedAt`, `UpdatedAt` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := gorm.NowFunc.Unix()
		if createTimeField, ok := scope.FieldByName("CreatedAt"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("UpdatedAt"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `UpdatedAt` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("UpdatedAt", time.Now().Unix())
	}
}

// deleteCallback will set `DeletedAt` where deleting
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedAt")

		if !scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

// addExtraSpaceIfExist adds a separator
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
