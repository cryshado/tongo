package abi

import (
	"context"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/tonkeeper/tongo/liteapi"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/tvm"
)

func mustToAddress(x string) tongo.AccountID {
	accountID, err := tongo.AccountIDFromRaw(x)
	if err != nil {
		panic(err)
	}
	return accountID
}

func mustToMsgAddress(x string) tlb.MsgAddress {
	accountID, err := tongo.AccountIDFromRaw(x)
	if err != nil {
		panic(err)
	}
	addr := tlb.MsgAddress{
		SumType: "AddrStd",
		AddrStd: struct {
			Anycast     tlb.Maybe[tlb.Anycast]
			WorkchainId int8
			Address     tlb.Bits256
		}{
			WorkchainId: int8(accountID.Workchain),
			Address:     accountID.Address,
		},
	}
	return addr
}

func TestGetMethods(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		data         string
		account      string
		method       InvokeFn
		want         any
		wantTypeHint string
	}{
		{
			name:         "GetSaleData for getgems auction",
			code:         "b5ee9c7201023501000b69000114ff00f4a413f4bcf2c80b010201200232020148032f0202ce042c020120052b04f5007434c0c05c6c3c903e900c36cf3e10b03ffe10d48831c16c23b40ccc74c7c87000234127265706561745f656e645f61756374696f6e8148831c16c23a0d6f6cf380cb00023411656d657267656e63795f6d6573736167658148831c16c26b50c3434c1f50c007ec0380c383e14d48431c163a10ccc76cf380060331f0c0e00caf8416edded44d0d20001f862d20001f864d20001f866fa4001f86dfa0001f86ed31f01f86fd31f01f870fa4001f872d401f868d430f869f849d0d21f01f867fa4001f863fa0001f86afa0001f86bfa0001f86cd31f01f871fa4001f873d31f30f8657ff8610294f84ec0008e3d7020f82582105fcc3d14c8cb1fcb3ff852cf165003cf1612cb0021fa02cb00c9718018c8cb05f853cf1670fa02cb6acc82080f424070fb02c98306fb00e30e7ff862db3c203402fadb3cf84e4054f00320c2008e2b70208010c8cb055007cf1622fa0216cb6a15cb1f8bf4d61726b6574706c616365206665658cf16c972fb009134e2f84e4003f00320c2008e2370208010c8cb055004cf1622fa0213cb6a12cb1f8b7526f79616c74798cf16c972fb009131e282080f424070fb02f84e58a101a120c20031220020f848d0fa40d31fd31ffa40d31fd31f3000c08e2270208010c8cb05f852cf165003fa0212cb6acb1f8b650726f6669748cf16c972fb009130e27020f82582105fcc3d14c8cb1fcb3ff84dcf165003cf1612cb008208989680fa02cb00c9718018c8cb05f853cf1670fa02cb6accc98306fb000054f849f848f850f84ff846f844f842c8ca00ca00ca00f84dcf16f84efa02cb1fcb1ff852cf16ccccc9ed54015c318103e9f852d749c202f2f28103ea01d31f821005138d9112ba12f2f48040d721fa4030f87270f8627ff864db3c340054f849f848f850f84ff846f844f842c8ca00ca00ca00f84dcf16f84efa02cb1fcb1ff852cf16ccccc9ed540486db3c20c0018f38308103edf823f850bef2f28103edf842c0fff2f28103f00282103b9aca00b912f2f2f8525210c705f8435220c705b1f2e193017fdb3cdb3ce020c0020f271315008c20c700c0ff923070e0d31f318b663616e63656c821c705923071e08b473746f70821c705923072e08b666696e697368821c705923072e08b66465706c6f79801c7059173e07002f2f84ec101915be0f84ef847a1228208989680a15210bc9930018208989680a1019132e28d0a565bdd5c88189a59081a185cc81899595b881bdd5d189a5908189e48185b9bdd1a195c881d5cd95c8ba001c0ff8e1f308d06d05d58dd1a5bdb881a185cc81899595b8818d85b98d95b1b19590ba0de21c200e30f2829003870208018c8cb05f84dcf165004fa0213cb6a12cb1f01cf16c972fb0000025b018a7020f82582105fcc3d14c8cb1fcb3ff852cf165003cf1612cb0021fa02cb00c9718018c8cb05f853cf1670fa02cb6acc82080f424070fb02c98306fb007ff8627ff866db3c340054f849f848f850f84ff846f844f842c8ca00ca00ca00f84dcf16f84efa02cb1fcb1ff852cf16ccccc9ed5404fc8ec330328103edf842c0fff2f28103f00182103b9aca00b9f2f2f823f850be8e17f8525210c705f8435220c705b1f84d5220c705b1f2e19399f8525210c705f2e193e2db3ce0c003925f03e0f842c0fff823f850beb1975f038103edf2f0e0f84b82103b9aca00a05220bef84bc200b0e302f850f851a1f823b9e300f84e1f1b24250294f84ec0008e3d7020f82582105fcc3d14c8cb1fcb3ff852cf165003cf1612cb0021fa02cb00c9718018c8cb05f853cf1670fa02cb6acc82080f424070fb02c98306fb00e30e7ff862db3c203402fadb3cf84e4054f00320c2008e2b70208010c8cb055007cf1622fa0216cb6a15cb1f8bf4d61726b6574706c616365206665658cf16c972fb009134e2f84e4003f00320c2008e2370208010c8cb055004cf1622fa0213cb6a12cb1f8b7526f79616c74798cf16c972fb009131e282080f424070fb02f84e58a101a120c20031220020f848d0fa40d31fd31ffa40d31fd31f3000c08e2270208010c8cb05f852cf165003fa0212cb6acb1f8b650726f6669748cf16c972fb009130e27020f82582105fcc3d14c8cb1fcb3ff84dcf165003cf1612cb008208989680fa02cb00c9718018c8cb05f853cf1670fa02cb6accc98306fb000054f849f848f850f84ff846f844f842c8ca00ca00ca00f84dcf16f84efa02cb1fcb1ff852cf16ccccc9ed54022c0270db3c21f86d82103b9aca00a1f86ef823f86fdb3c271f02f2f84ec101915be0f84ef847a1228208989680a15210bc9930018208989680a1019132e28d0a565bdd5c88189a59081a185cc81899595b881bdd5d189a5908189e48185b9bdd1a195c881d5cd95c8ba001c0ff8e1f308d06d05d58dd1a5bdb881a185cc81899595b8818d85b98d95b1b19590ba0de21c200e30f2829003870208018c8cb05f84dcf165004fa0213cb6a12cb1f01cf16c972fb0000025b0294f84ec0008e3d7020f82582105fcc3d14c8cb1fcb3ff852cf165003cf1612cb0021fa02cb00c9718018c8cb05f853cf1670fa02cb6acc82080f424070fb02c98306fb00e30e7ff862db3c203402fadb3cf84e4054f00320c2008e2b70208010c8cb055007cf1622fa0216cb6a15cb1f8bf4d61726b6574706c616365206665658cf16c972fb009134e2f84e4003f00320c2008e2370208010c8cb055004cf1622fa0213cb6a12cb1f8b7526f79616c74798cf16c972fb009131e282080f424070fb02f84e58a101a120c20031220020f848d0fa40d31fd31ffa40d31fd31f3000c08e2270208010c8cb05f852cf165003fa0212cb6acb1f8b650726f6669748cf16c972fb009130e27020f82582105fcc3d14c8cb1fcb3ff84dcf165003cf1612cb008208989680fa02cb00c9718018c8cb05f853cf1670fa02cb6accc98306fb000054f849f848f850f84ff846f844f842c8ca00ca00ca00f84dcf16f84efa02cb1fcb1ff852cf16ccccc9ed54000ef850f851a0f87003708e95328103e8f84a5220b9f2f2f86ef86df823f86fdb3ce1f84ef84ca05220b9975f038103e8f2f0e00270db3c01f86df86ef823f86fdb3c3427340054f849f848f850f84ff846f844f842c8ca00ca00ca00f84dcf16f84efa02cb1fcb1ff852cf16ccccc9ed5402f2f84ec101915be0f84ef847a1228208989680a15210bc9930018208989680a1019132e28d0a565bdd5c88189a59081a185cc81899595b881bdd5d189a5908189e48185b9bdd1a195c881d5cd95c8ba001c0ff8e1f308d06d05d58dd1a5bdb881a185cc81899595b8818d85b98d95b1b19590ba0de21c200e30f2829003870208018c8cb05f84dcf165004fa0213cb6a12cb1f01cf16c972fb0000025b0054f849f848f850f84ff846f844f842c8ca00ca00ca00f84dcf16f84efa02cb1fcb1ff852cf16ccccc9ed54001320840ee6b280006a61200201202d2e001120840ee6b2802a6120001d08300024d7c0dc38167c00807c0060028ba03859b679b679041082aa87f085f0a1f087f0a7f0a5f09df09bf099f097f095f08bf09ff08c1a22261a182224181622221614222014213e211c20fa20d820b620f420d220b1333100caf8416edded44d0d20001f862d20001f864d20001f866fa4001f86dfa0001f86ed31f01f86fd31f01f870fa4001f872d401f868d430f869f849d0d21f01f867fa4001f863fa0001f86afa0001f86bfa0001f86cd31f01f871fa4001f873d31f30f8657ff8610020f848d0fa40d31fd31ffa40d31fd31f300228f230db3c8103eef844c0fff2f2f8007ff864db3c333400caf8416edded44d0d20001f862d20001f864d20001f866fa4001f86dfa0001f86ed31f01f86fd31f01f870fa4001f872d401f868d430f869f849d0d21f01f867fa4001f863fa0001f86afa0001f86bfa0001f86cd31f01f871fa4001f873d31f30f8657ff8610054f849f848f850f84ff846f844f842c8ca00ca00ca00f84dcf16f84efa02cb1fcb1ff852cf16ccccc9ed54",
			data:         "b5ee9c720101030100e0000255400000000031fcf510c00047efbb81020160d28464124271660d1565daad67fa6f984a4f4c116586f59688010200a58014726b0c3ef3b5eb3427ada305c2c80421805f31c7be31fb4e971eb5628357e300000000a000000c900182cec523a039304feb0fd0fbf8df6036f9f5ca825269b47351699c4878be7b08000000280000019200b30080ebe8800b09dcc365bfe106e22da1f96a0f1b272c9797d380bfad4283637f94bad487c30a020c855800087735940000000259001270316eb10f404fd1c575e336fb68e0a8d82a605c7d89a66b7a5a054883ac5758fdd1c8a0",
			account:      "0:6c912c30dc3642d4f4add23c24485312aa351a804fe5294dfca11f83e849d10d",
			method:       GetSaleData,
			wantTypeHint: "GetSaleData_GetgemsAuctionResult",
			want: GetSaleData_GetgemsAuctionResult{
				Magic:            4281667,
				End:              false,
				EndTime:          1677322785,
				Marketplace:      mustToMsgAddress("0:584ee61b2dff0837116d0fcb5078d93964bcbe9c05fd6a141b1bfca5d6a43e18"),
				Nft:              mustToMsgAddress("0:49c0c5bac43d013f4715d78cdbeda382a360a98171f62699ade96815220eb15d"),
				Owner:            mustToMsgAddress("0:047efbb81020160d28464124271660d1565daad67fa6f984a4f4c116586f5968"),
				LastBid:          0,
				LastMember:       tlb.MsgAddress{SumType: "AddrNone"},
				MinStep:          1000000000,
				MarketFeeAddress: mustToMsgAddress("0:a3935861f79daf59a13d6d182e1640210c02f98e3df18fda74b8f5ab141abf18"),
				MpFeeFactor:      5,
				MpFeeBase:        100,
				RoyaltyAddress:   mustToMsgAddress("0:60b3b148e80e4c13fac3f43efe37d80dbe7d72a0949a6d1cd45a67121e2f9ec2"),
				RoyaltyFeeFactor: 10,
				RoyaltyFeeBase:   100,
				MaxBid:           0,
				MinBid:           4400000000,
				CreatedAt:        1677149986,
				LastBidAt:        0,
				IsCanceled:       false,
			},
		},
		{
			name:         "GetPluginList",
			code:         "b5ee9c72010214010002d4000114ff00f4a413f4bcf2c80b01020120020f020148030602e6d001d0d3032171b0925f04e022d749c120925f04e002d31f218210706c7567bd22821064737472bdb0925f05e003fa403020fa4401c8ca07cbffc9d0ed44d0810140d721f404305c810108f40a6fa131b3925f07e005d33fc8258210706c7567ba923830e30d03821064737472ba925f06e30d0405007801fa00f40430f8276f2230500aa121bef2e0508210706c7567831eb17080185004cb0526cf1658fa0219f400cb6917cb1f5260cb3f20c98040fb0006008a5004810108f45930ed44d0810140d720c801cf16f400c9ed540172b08e23821064737472831eb17080185005cb055003cf1623fa0213cb6acb1fcb3fc98040fb00925f03e2020120070e020120080d020158090a003db29dfb513420405035c87d010c00b23281f2fff274006040423d029be84c600201200b0c0019adce76a26840206b90eb85ffc00019af1df6a26840106b90eb858fc00011b8c97ed44d0d70b1f80059bd242b6f6a2684080a06b90fa0218470d4080847a4937d29910ce6903e9ff9837812801b7810148987159f318404f8f28308d71820d31fd31fd31f02f823bbf264ed44d0d31fd31fd3fff404d15143baf2a15151baf2a205f901541064f910f2a3f80024a4c8cb1f5240cb1f5230cbff5210f400c9ed54f80f01d30721c0009f6c519320d74a96d307d402fb00e830e021c001e30021c002e30001c0039130e30d03a4c8cb1f12cb1fcbff10111213006ed207fa00d4d422f90005c8ca0715cbffc9d077748018c8cb05cb0222cf165005fa0214cb6b12ccccc973fb00c84014810108f451f2a7020070810108d718fa00d33fc8542047810108f451f2a782106e6f746570748018c8cb05cb025006cf165004fa0214cb6a12cb1fcb3fc973fb0002006c810108d718fa00d33f305224810108f459f2a782106473747270748018c8cb05cb025005cf165003fa0213cb6acb1f12cb3fc973fb00000af400c9ed54",
			data:         "b5ee9c720101060100a100015100001e8f29a9a3179190b2045703bc50371dd6897707b1791b8c4cbdd775e41b1ab7b3ebc8ae467bc001020581400c020502027303040041be5d7739ed1b853cab49680427bd0c78a375531c7a203092a115f1f9033cb7f6d00041be4d31620ef663745b654fb7bd023166f67fc8097fbfb34a5f6d60eb4480a303900041bf618db178d663998eaf849ac43cde4965aeac444373c92a65a10f10a2e6f406af",
			account:      "0:9a7752cab755c829967b33e7f2692f9bdb81a47415168bb89c10b74ee0defc6b",
			method:       GetPluginList,
			wantTypeHint: "GetPluginListResult",
			want: GetPluginListResult{
				Plugins: []struct {
					Workchain int32
					Address   tlb.Bits256
				}{
					{Workchain: 0, Address: mustToAddress("0:70c6d8bc6b31ccc757c24d621e6f24b2d7562221b9e49532d0878851737a0357").Address},
					{Workchain: 0, Address: mustToAddress("0:4e698b1077b31ba2db2a7dbde8118b37b3fe404bfdfd9a52fb6b075a2405181c").Address},
					{Workchain: 0, Address: mustToAddress("0:4cebb9cf68dc29e55a4b40213de863c51baa98e3d101849508af8fc819e5bfb6").Address},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mainnetConfig, _ := boc.DeserializeBocBase64(mainnetConfig)
			code, err := hex.DecodeString(tt.code)
			if err != nil {
				t.Fatalf("DecodeString() failed: %v", err)
			}
			data, err := hex.DecodeString(tt.data)
			if err != nil {
				t.Fatalf("DecodeString() failed: %v", err)
			}
			accountID, err := tongo.AccountIDFromRaw(tt.account)
			if err != nil {
				t.Fatalf("AccountIDFromRaw() failed: %v", err)
			}
			codeCell, _ := boc.DeserializeBoc(code)
			dataCell, _ := boc.DeserializeBoc(data)
			emulator, err := tvm.NewEmulator(codeCell[0], dataCell[0], mainnetConfig[0], 1_000_000_000, -1)
			typeHint, got, err := tt.method(context.Background(), emulator, accountID)
			if err != nil {
				t.Fatalf("GetPluginList() failed: %v", err)
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("want: %v, got: %v", tt.want, got)
			}
			if !reflect.DeepEqual(tt.wantTypeHint, typeHint) {
				t.Fatalf("want: %v, got: %v", tt.wantTypeHint, typeHint)
			}
		})
	}
}

func TestWhalesNominators(t *testing.T) {
	address := tongo.MustParseAccountID("EQBI-wGVp_x0VFEjd7m9cEUD3tJ_bnxMSp0Tb9qz757ATEAM")
	client, err := liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		t.Fatal(err)
	}
	_, v, err := GetMembers(context.Background(), client, address)
	if err != nil {
		t.Fatal(err)
	}
	members := v.(GetMembers_WhalesNominatorResult).Members
	if len(members) == 0 || members[0].Address.SumType != "AddrStd" {
		t.Fatal(len(members))
	}
	_, v, err = GetPoolStatus(context.Background(), client, address)
	if err != nil {
		t.Fatal(err)
	}
	status := v.(GetPoolStatusResult)
	fmt.Printf("%+v\n", status)

	_, v, err = GetParams(context.Background(), client, address)
	if err != nil {
		t.Fatal(err)
	}
	params := v.(GetParams_WhalesNominatorResult)
	fmt.Printf("%+v\n", params)

}

func TestMethodsDecode(t *testing.T) {
	c, err := boc.DeserializeSinglRootBase64("te6ccgEBBAEA7AADs0Y3KJoZbSToF2FWrsk5n2kJkqyX4X6Ap8VP92juX4NPlNJSZUBwdJYp6pn3SVlg0xt+7QjJLdBJYx7JVdtEr9ZqVVgPAAAAA2QEUxhkBFOuCHpoZW5nc2h1wAECAwBgAWh0dHBzOi8vbmZ0LmZyYWdtZW50LmNvbS91c2VybmFtZS96aGVuZ3NodS5qc29uAGGACBG0dlFtgMtLJ8IHIk03VVOL8ZXXgY07PhVXdWVTJK1qMG3EIAAAoAABwgABJ1AQAEsABQBkgAgRtHZRbYDLSyfCByJNN1VTi/GV14GNOz4VV3VlUyStcA==")
	if err != nil {
		return
	}
	typeName, v, err := MessageDecoder(c)
	if err != nil {
		t.Fatal(err)
	}
	if typeName != "TelemintDeploy" {
		t.Fatal(typeName)
	}
	body := v.(TelemintDeployMsgBody)
	if body.Msg.Username != "zhengshu" {
		t.Fatal(body.Msg.Username)
	}

}
