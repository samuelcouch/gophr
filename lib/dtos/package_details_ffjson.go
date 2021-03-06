// DO NOT EDIT!
// Code generated by ffjson <https://github.com/pquerna/ffjson>
// source: package_details.go
// DO NOT EDIT!

package dtos

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	fflib "github.com/pquerna/ffjson/fflib/v1"
)

func (mj *PackageDetails) MarshalJSON() ([]byte, error) {
	var buf fflib.Buffer
	if mj == nil {
		buf.WriteString("null")
		return buf.Bytes(), nil
	}
	err := mj.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (mj *PackageDetails) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if mj == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	buf.WriteString(`{"repo":`)
	fflib.WriteJsonString(buf, string(mj.Repo))
	buf.WriteString(`,"stars":`)
	fflib.FormatBits2(buf, uint64(mj.Stars), 10, mj.Stars < 0)
	buf.WriteString(`,"author":`)
	fflib.WriteJsonString(buf, string(mj.Author))
	if mj.Awesome {
		buf.WriteString(`,"awesome":true`)
	} else {
		buf.WriteString(`,"awesome":false`)
	}
	buf.WriteString(`,"versions":`)
	if mj.Versions != nil {
		buf.WriteString(`[`)
		for i, v := range mj.Versions {
			if i != 0 {
				buf.WriteString(`,`)
			}
			/* Struct fall back. type=dtos.PackageVersion kind=struct */
			err = buf.Encode(&v)
			if err != nil {
				return err
			}
		}
		buf.WriteString(`]`)
	} else {
		buf.WriteString(`null`)
	}
	/* Struct fall back. type=dtos.PackageDownloads kind=struct */
	buf.WriteString(`,"downloads":`)
	err = buf.Encode(&mj.Downloads)
	if err != nil {
		return err
	}
	buf.WriteString(`,"trendScore":`)
	fflib.AppendFloat(buf, float64(mj.TrendScore), 'g', -1, 32)
	buf.WriteString(`,"description":`)
	fflib.WriteJsonString(buf, string(mj.Description))
	buf.WriteString(`,"dateDiscovered":`)

	{

		obj, err = mj.DateDiscovered.MarshalJSON()
		if err != nil {
			return err
		}
		buf.Write(obj)

	}
	buf.WriteString(`,"dateLastIndexed":`)

	{

		obj, err = mj.DateLastIndexed.MarshalJSON()
		if err != nil {
			return err
		}
		buf.Write(obj)

	}
	buf.WriteByte('}')
	return nil
}

const (
	ffj_t_PackageDetailsbase = iota
	ffj_t_PackageDetailsno_such_key

	ffj_t_PackageDetails_Repo

	ffj_t_PackageDetails_Stars

	ffj_t_PackageDetails_Author

	ffj_t_PackageDetails_Awesome

	ffj_t_PackageDetails_Versions

	ffj_t_PackageDetails_Downloads

	ffj_t_PackageDetails_TrendScore

	ffj_t_PackageDetails_Description

	ffj_t_PackageDetails_DateDiscovered

	ffj_t_PackageDetails_DateLastIndexed
)

var ffj_key_PackageDetails_Repo = []byte("repo")

var ffj_key_PackageDetails_Stars = []byte("stars")

var ffj_key_PackageDetails_Author = []byte("author")

var ffj_key_PackageDetails_Awesome = []byte("awesome")

var ffj_key_PackageDetails_Versions = []byte("versions")

var ffj_key_PackageDetails_Downloads = []byte("downloads")

var ffj_key_PackageDetails_TrendScore = []byte("trendScore")

var ffj_key_PackageDetails_Description = []byte("description")

var ffj_key_PackageDetails_DateDiscovered = []byte("dateDiscovered")

var ffj_key_PackageDetails_DateLastIndexed = []byte("dateLastIndexed")

func (uj *PackageDetails) UnmarshalJSON(input []byte) error {
	fs := fflib.NewFFLexer(input)
	return uj.UnmarshalJSONFFLexer(fs, fflib.FFParse_map_start)
}

func (uj *PackageDetails) UnmarshalJSONFFLexer(fs *fflib.FFLexer, state fflib.FFParseState) error {
	var err error = nil
	currentKey := ffj_t_PackageDetailsbase
	_ = currentKey
	tok := fflib.FFTok_init
	wantedTok := fflib.FFTok_init

mainparse:
	for {
		tok = fs.Scan()
		//	println(fmt.Sprintf("debug: tok: %v  state: %v", tok, state))
		if tok == fflib.FFTok_error {
			goto tokerror
		}

		switch state {

		case fflib.FFParse_map_start:
			if tok != fflib.FFTok_left_bracket {
				wantedTok = fflib.FFTok_left_bracket
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_key
			continue

		case fflib.FFParse_after_value:
			if tok == fflib.FFTok_comma {
				state = fflib.FFParse_want_key
			} else if tok == fflib.FFTok_right_bracket {
				goto done
			} else {
				wantedTok = fflib.FFTok_comma
				goto wrongtokenerror
			}

		case fflib.FFParse_want_key:
			// json {} ended. goto exit. woo.
			if tok == fflib.FFTok_right_bracket {
				goto done
			}
			if tok != fflib.FFTok_string {
				wantedTok = fflib.FFTok_string
				goto wrongtokenerror
			}

			kn := fs.Output.Bytes()
			if len(kn) <= 0 {
				// "" case. hrm.
				currentKey = ffj_t_PackageDetailsno_such_key
				state = fflib.FFParse_want_colon
				goto mainparse
			} else {
				switch kn[0] {

				case 'a':

					if bytes.Equal(ffj_key_PackageDetails_Author, kn) {
						currentKey = ffj_t_PackageDetails_Author
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffj_key_PackageDetails_Awesome, kn) {
						currentKey = ffj_t_PackageDetails_Awesome
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'd':

					if bytes.Equal(ffj_key_PackageDetails_Downloads, kn) {
						currentKey = ffj_t_PackageDetails_Downloads
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffj_key_PackageDetails_Description, kn) {
						currentKey = ffj_t_PackageDetails_Description
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffj_key_PackageDetails_DateDiscovered, kn) {
						currentKey = ffj_t_PackageDetails_DateDiscovered
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffj_key_PackageDetails_DateLastIndexed, kn) {
						currentKey = ffj_t_PackageDetails_DateLastIndexed
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'r':

					if bytes.Equal(ffj_key_PackageDetails_Repo, kn) {
						currentKey = ffj_t_PackageDetails_Repo
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 's':

					if bytes.Equal(ffj_key_PackageDetails_Stars, kn) {
						currentKey = ffj_t_PackageDetails_Stars
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 't':

					if bytes.Equal(ffj_key_PackageDetails_TrendScore, kn) {
						currentKey = ffj_t_PackageDetails_TrendScore
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'v':

					if bytes.Equal(ffj_key_PackageDetails_Versions, kn) {
						currentKey = ffj_t_PackageDetails_Versions
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				}

				if fflib.EqualFoldRight(ffj_key_PackageDetails_DateLastIndexed, kn) {
					currentKey = ffj_t_PackageDetails_DateLastIndexed
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffj_key_PackageDetails_DateDiscovered, kn) {
					currentKey = ffj_t_PackageDetails_DateDiscovered
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffj_key_PackageDetails_Description, kn) {
					currentKey = ffj_t_PackageDetails_Description
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffj_key_PackageDetails_TrendScore, kn) {
					currentKey = ffj_t_PackageDetails_TrendScore
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffj_key_PackageDetails_Downloads, kn) {
					currentKey = ffj_t_PackageDetails_Downloads
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffj_key_PackageDetails_Versions, kn) {
					currentKey = ffj_t_PackageDetails_Versions
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffj_key_PackageDetails_Awesome, kn) {
					currentKey = ffj_t_PackageDetails_Awesome
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffj_key_PackageDetails_Author, kn) {
					currentKey = ffj_t_PackageDetails_Author
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffj_key_PackageDetails_Stars, kn) {
					currentKey = ffj_t_PackageDetails_Stars
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffj_key_PackageDetails_Repo, kn) {
					currentKey = ffj_t_PackageDetails_Repo
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				currentKey = ffj_t_PackageDetailsno_such_key
				state = fflib.FFParse_want_colon
				goto mainparse
			}

		case fflib.FFParse_want_colon:
			if tok != fflib.FFTok_colon {
				wantedTok = fflib.FFTok_colon
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_value
			continue
		case fflib.FFParse_want_value:

			if tok == fflib.FFTok_left_brace || tok == fflib.FFTok_left_bracket || tok == fflib.FFTok_integer || tok == fflib.FFTok_double || tok == fflib.FFTok_string || tok == fflib.FFTok_bool || tok == fflib.FFTok_null {
				switch currentKey {

				case ffj_t_PackageDetails_Repo:
					goto handle_Repo

				case ffj_t_PackageDetails_Stars:
					goto handle_Stars

				case ffj_t_PackageDetails_Author:
					goto handle_Author

				case ffj_t_PackageDetails_Awesome:
					goto handle_Awesome

				case ffj_t_PackageDetails_Versions:
					goto handle_Versions

				case ffj_t_PackageDetails_Downloads:
					goto handle_Downloads

				case ffj_t_PackageDetails_TrendScore:
					goto handle_TrendScore

				case ffj_t_PackageDetails_Description:
					goto handle_Description

				case ffj_t_PackageDetails_DateDiscovered:
					goto handle_DateDiscovered

				case ffj_t_PackageDetails_DateLastIndexed:
					goto handle_DateLastIndexed

				case ffj_t_PackageDetailsno_such_key:
					err = fs.SkipField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}
					state = fflib.FFParse_after_value
					goto mainparse
				}
			} else {
				goto wantedvalue
			}
		}
	}

handle_Repo:

	/* handler: uj.Repo type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			uj.Repo = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Stars:

	/* handler: uj.Stars type=int kind=int quoted=false*/

	{
		if tok != fflib.FFTok_integer && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for int", tok))
		}
	}

	{

		if tok == fflib.FFTok_null {

		} else {

			tval, err := fflib.ParseInt(fs.Output.Bytes(), 10, 64)

			if err != nil {
				return fs.WrapErr(err)
			}

			uj.Stars = int(tval)

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Author:

	/* handler: uj.Author type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			uj.Author = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Awesome:

	/* handler: uj.Awesome type=bool kind=bool quoted=false*/

	{
		if tok != fflib.FFTok_bool && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for bool", tok))
		}
	}

	{
		if tok == fflib.FFTok_null {

		} else {
			tmpb := fs.Output.Bytes()

			if bytes.Compare([]byte{'t', 'r', 'u', 'e'}, tmpb) == 0 {

				uj.Awesome = true

			} else if bytes.Compare([]byte{'f', 'a', 'l', 's', 'e'}, tmpb) == 0 {

				uj.Awesome = false

			} else {
				err = errors.New("unexpected bytes for true/false value")
				return fs.WrapErr(err)
			}

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Versions:

	/* handler: uj.Versions type=[]dtos.PackageVersion kind=slice quoted=false*/

	{

		{
			if tok != fflib.FFTok_left_brace && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for ", tok))
			}
		}

		if tok == fflib.FFTok_null {
			uj.Versions = nil
		} else {

			uj.Versions = []PackageVersion{}

			wantVal := true

			for {

				var tmp_uj__Versions PackageVersion

				tok = fs.Scan()
				if tok == fflib.FFTok_error {
					goto tokerror
				}
				if tok == fflib.FFTok_right_brace {
					break
				}

				if tok == fflib.FFTok_comma {
					if wantVal == true {
						// TODO(pquerna): this isn't an ideal error message, this handles
						// things like [,,,] as an array value.
						return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
					}
					continue
				} else {
					wantVal = true
				}

				/* handler: tmp_uj__Versions type=dtos.PackageVersion kind=struct quoted=false*/

				{
					/* Falling back. type=dtos.PackageVersion kind=struct */
					tbuf, err := fs.CaptureField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}

					err = json.Unmarshal(tbuf, &tmp_uj__Versions)
					if err != nil {
						return fs.WrapErr(err)
					}
				}

				uj.Versions = append(uj.Versions, tmp_uj__Versions)

				wantVal = false
			}
		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Downloads:

	/* handler: uj.Downloads type=dtos.PackageDownloads kind=struct quoted=false*/

	{
		/* Falling back. type=dtos.PackageDownloads kind=struct */
		tbuf, err := fs.CaptureField(tok)
		if err != nil {
			return fs.WrapErr(err)
		}

		err = json.Unmarshal(tbuf, &uj.Downloads)
		if err != nil {
			return fs.WrapErr(err)
		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_TrendScore:

	/* handler: uj.TrendScore type=float32 kind=float32 quoted=false*/

	{
		if tok != fflib.FFTok_double && tok != fflib.FFTok_integer && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for float32", tok))
		}
	}

	{

		if tok == fflib.FFTok_null {

		} else {

			tval, err := fflib.ParseFloat(fs.Output.Bytes(), 32)

			if err != nil {
				return fs.WrapErr(err)
			}

			uj.TrendScore = float32(tval)

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Description:

	/* handler: uj.Description type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			uj.Description = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_DateDiscovered:

	/* handler: uj.DateDiscovered type=time.Time kind=struct quoted=false*/

	{
		if tok == fflib.FFTok_null {

			state = fflib.FFParse_after_value
			goto mainparse
		}

		tbuf, err := fs.CaptureField(tok)
		if err != nil {
			return fs.WrapErr(err)
		}

		err = uj.DateDiscovered.UnmarshalJSON(tbuf)
		if err != nil {
			return fs.WrapErr(err)
		}
		state = fflib.FFParse_after_value
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_DateLastIndexed:

	/* handler: uj.DateLastIndexed type=time.Time kind=struct quoted=false*/

	{
		if tok == fflib.FFTok_null {

			state = fflib.FFParse_after_value
			goto mainparse
		}

		tbuf, err := fs.CaptureField(tok)
		if err != nil {
			return fs.WrapErr(err)
		}

		err = uj.DateLastIndexed.UnmarshalJSON(tbuf)
		if err != nil {
			return fs.WrapErr(err)
		}
		state = fflib.FFParse_after_value
	}

	state = fflib.FFParse_after_value
	goto mainparse

wantedvalue:
	return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
wrongtokenerror:
	return fs.WrapErr(fmt.Errorf("ffjson: wanted token: %v, but got token: %v output=%s", wantedTok, tok, fs.Output.String()))
tokerror:
	if fs.BigError != nil {
		return fs.WrapErr(fs.BigError)
	}
	err = fs.Error.ToError()
	if err != nil {
		return fs.WrapErr(err)
	}
	panic("ffjson-generated: unreachable, please report bug.")
done:
	return nil
}
