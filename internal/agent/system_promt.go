package agent

import (
	"os"      // [THEM_MOI]
	"strings" // [THEM_MOI]
)

func (a *Agent) SystemPrompt() string {
	sys := `[VAI TRO]: Nhan vien phuc vu AI chuyen nghiep.
[THONG TIN QUAN]: {info}
[MENU]: {menu}
[BOI CANH KHACH HANG]:
- Khach hang: {customer_name}
- SDT: {customer_phone}
- Ban so: {table}
- Don hang hien tai: {current_items}
[NHIEM VU]: Tu van va chot don. Luon tra ve JSON theo Schema da dinh nghia, chi gom text, photourl, videourl. // [SUA]
[TOOL create_order]: Bat buoc goi function create_order khi khach da xac nhan dat hang va da du thong tin bat buoc. // [SUA]
- Thong tin bat buoc: danh sach mon co name va quantity, dining_option.
- Neu dining_option la "tai cho" thi phai co table_number.
- Neu dining_option la "mang ve" thi phai co customer.name hoac customer.phone.
- Neu thieu thong tin, khong goi function va text phai hoi dung thong tin con thieu. // [SUA]
- Khi du thong tin va khach xac nhan, goi function create_order voi day du items, dining_option, customer, table_number neu co, note neu co. // [SUA]
- Sau khi he thong tra function_result, tra JSON cuoi cung chi gom text, photourl, videourl. // [SUA]` // [SUA]

	menuBytes, err := os.ReadFile("data/menu.txt") // [THEM_MOI]
	if err == nil {
		sys = strings.ReplaceAll(sys, "{menu}", string(menuBytes)) // [THEM_MOI]
	} else {
		sys = strings.ReplaceAll(sys, "{menu}", "") // [THEM_MOI]
	}

	infoBytes, err := os.ReadFile("data/info.txt") // [THEM_MOI]
	if err != nil {
		infoBytes, err = os.ReadFile("data/infor.txt") // [THEM_MOI]
	}
	if err == nil {
		sys = strings.ReplaceAll(sys, "{info}", string(infoBytes)) // [THEM_MOI]
	} else {
		sys = strings.ReplaceAll(sys, "{info}", "") // [THEM_MOI]
	}

	return sys
}
