package api

import (
	"testing"
)

func TestNewProductCode(t *testing.T) {
	type args struct {
		productCode string
	}
	tests := []struct {
		name    string
		args    args
		want    ProductCode
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				productCode: ProductCodeBTCJPY,
			},
			want:    ProductCode(ProductCodeBTCJPY),
			wantErr: false,
		},
		{
			name: "validate error",
			args: args{
				productCode: "invalide",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProductCode(tt.args.productCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProductCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewProductCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductCode_Validate(t *testing.T) {
	tests := []struct {
		name string
		p    ProductCode
		want bool
	}{
		{
			name: "success",
			p:    ProductCode(ProductCodeBTCJPY),
			want: true,
		},
		{
			name: "invalid",
			p:    ProductCode("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Validate(); got != tt.want {
				t.Errorf("ProductCode.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
