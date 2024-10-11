// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package body

type Options struct {
	ApnsProduction    bool              `json:"apns_production"`
	ThirdPartyChannel ThirdPartyChannel `json:"third_party_channel"`
}

type ThirdPartyChannel struct {
	Vivo   Vivo   `json:"vivo,omitempty"`
	Xiaomi Xiaomi `json:"xiaomi,omitempty"`
}

type Vivo struct {
	Classification int `json:"classification"`
}

type Xiaomi struct {
	Channel_id string `json:"channel_id"`
	Skip_quota bool   `json:"skip_quota"`
}

func (o *Options) SetApnsProduction(c bool) {
	o.ApnsProduction = c
}

func (o *Options) SetThirdPartyChannel() {
	o.ThirdPartyChannel.Vivo.Classification = 1
	o.ThirdPartyChannel.Xiaomi.Channel_id = "127711"
	o.ThirdPartyChannel.Xiaomi.Skip_quota = true
}
