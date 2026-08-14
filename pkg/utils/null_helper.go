package utils

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/sqlc-dev/pqtype"
)

// NullStringToFloatPointer แปลงเลขที่อยู่ในรูป sql.NullString ให้กลายเป็น *float64
func NullStringToFloatPointer(ns sql.NullString) *float64 {
	// 1. ถ้าค่าเป็น NULL ใน Database ให้รีเทิร์น nil
	if !ns.Valid {
		return nil
	}

	// 2. พยายามแปลง String เป็น Float64
	val, err := strconv.ParseFloat(ns.String, 64)
	if err != nil {
		// ถ้าข้อมูลใน DB พัง (เช่น เผลอเก็บตัวอักษรลงไป) จะคืนค่า nil ป้องกันระบบแครช
		return nil
	}

	// 3. ถ้าแปลงสำเร็จ ส่ง Pointer กลับไป
	return &val
}

// NullStringToPointer แปลง sql.NullString เป็น *string
func NullStringToPointer(ns sql.NullString) *string {
	if ns.Valid {
		// ต้องก๊อปปี้ค่าใส่ตัวแปรใหม่ก่อนรีเทิร์น เพื่อป้องกันบั๊กเรื่อง Pointer ในหน่วยความจำ
		val := ns.String
		return &val
	}
	return nil
}

// NullTimeToPointer แปลง sql.NullTime เป็น *time.Time
func NullTimeToPointer(nt sql.NullTime) *time.Time {
	if nt.Valid {
		val := nt.Time
		return &val
	}
	return nil
}

// NullFloat64ToPointer แปลง sql.NullFloat64 เป็น *float64
func NullFloat64ToPointer(nf sql.NullFloat64) *float64 {
	if nf.Valid {
		val := nf.Float64
		return &val
	}
	return nil
}

// NullInt32ToPointer แปลง sql.NullInt32 เป็น *int32
func NullInt32ToPointer(ni sql.NullInt32) *int32 {
	if ni.Valid {
		val := ni.Int32
		return &val
	}
	return nil
}

// NullBoolToPointer แปลง sql.NullBool เป็น *bool
func NullBoolToPointer(nb sql.NullBool) *bool {
	if nb.Valid {
		val := nb.Bool
		return &val
	}
	return nil
}

// PointerToNullString แปลง *string เป็น sql.NullString
func PointerToNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{Valid: false}
}

// PointerToNullBool แปลง *bool เป็น sql.NullBool
func PointerToNullBool(b *bool) sql.NullBool {
	if b != nil {
		return sql.NullBool{Bool: *b, Valid: true}
	}
	return sql.NullBool{Valid: false}
}

// PointerToNullFloat64 แปลง *float64 เป็น sql.NullFloat64
func PointerToNullFloat64(f *float64) sql.NullFloat64 {
	if f != nil {
		return sql.NullFloat64{Float64: *f, Valid: true}
	}
	return sql.NullFloat64{Valid: false}
}

// PointerToNullInt32 แปลง *int32 เป็น sql.NullInt32
func PointerToNullInt32(i *int32) sql.NullInt32 {
	if i != nil {
		return sql.NullInt32{Int32: *i, Valid: true}
	}
	return sql.NullInt32{Valid: false}
}

// PointerToNullTime แปลง *time.Time เป็น sql.NullTime
func PointerToNullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{Valid: false}
}

// 🟢 โบนัส: สำหรับจัดการ CustomizationFields ที่เป็น JSONB โดยเฉพาะ
// JsonToNullRawMessage แปลง json.RawMessage เป็น pqtype.NullRawMessage
func JsonToNullRawMessage(j json.RawMessage) pqtype.NullRawMessage {
	// เช็คว่ามีข้อมูลไหม และไม่ใช่ค่า "null" (แบบ String ที่ลอยมาจาก JSON)
	if len(j) > 0 && string(j) != "null" {
		return pqtype.NullRawMessage{RawMessage: j, Valid: true}
	}
	return pqtype.NullRawMessage{Valid: false}
}
