package paypal

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dghubble/sling"
)

func (c *Client) SetExpressCheckout(ctx context.Context, req SetExpressCheckoutReq) (SetExpressCheckoutResp, error) {
	path := c.config.Host + nvp

	request, err := sling.New().Post(path).BodyForm(setExpressCheckoutReq{
		NVPBase:               c.NvpBase(SetExpressCheckout),
		SetExpressCheckoutReq: req,
	}).Request()


	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return SetExpressCheckoutResp{}, err
	}

	defer resp.Body.Close()

	by ,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return SetExpressCheckoutResp{}, err
	}

	fmt.Println(string(by))

	value , err := url.ParseQuery(string(by))
	if err != nil {
		fmt.Println(err)
		return SetExpressCheckoutResp{}, err
	}

	return SetExpressCheckoutResp{Body: value}, nil
}

func (c *Client) NvpBase(method NVPMethod) NVPBase {
	return NVPBase{
		Method:    string(method),
		Version:   string(Version124),
		User:      c.config.NVPAndSOAPAPICredentials.Username,
		Pwd:       c.config.NVPAndSOAPAPICredentials.Password,
		Signature: c.config.NVPAndSOAPAPICredentials.Signature,
	}
}


func (c *Client) GetExpressCheckoutDetails(ctx context.Context, req GetExpressCheckoutDetailsReq) (GetExpressCheckoutDetailsResp, error) {
	path := c.config.Host + nvp

	request, err := sling.New().Post(path).BodyForm(getExpressCheckoutDetailsReq{
		NVPBase:               c.NvpBase(GetExpressCheckoutDetails),
		GetExpressCheckoutDetailsReq: req,
	}).Request()


	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return GetExpressCheckoutDetailsResp{}, err
	}

	defer resp.Body.Close()

	by ,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return GetExpressCheckoutDetailsResp{}, err
	}

	fmt.Println(string(by))

	value , err := url.ParseQuery(string(by))
	if err != nil {
		fmt.Println(err)
		return GetExpressCheckoutDetailsResp{}, err
	}

	return GetExpressCheckoutDetailsResp{Body: value}, nil
}

func (c *Client) DoExpressCheckoutPayment(ctx context.Context, req GetExpressCheckoutDetailsReq) (GetExpressCheckoutDetailsResp, error) {
	res ,err := c.GetExpressCheckoutDetails(ctx, req)
	if err != nil {
		fmt.Println(err)
		return GetExpressCheckoutDetailsResp{},err
	}

	form := setFormForDoExpressCheckoutPayment(c.config.NVPAndSOAPAPICredentials, res.Body,req.Token)

	client :=http.Client{}
	path := c.config.Host + nvp
	resp, err := client.PostForm(path, form)

	defer resp.Body.Close()

	by ,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return GetExpressCheckoutDetailsResp{}, err
	}

	fmt.Println(string(by))

	value , err := url.ParseQuery(string(by))
	if err != nil {
		fmt.Println(err)
		return GetExpressCheckoutDetailsResp{}, err
	}

	return GetExpressCheckoutDetailsResp{Body: value}, nil
}

func setFormForDoExpressCheckoutPayment(credentials APICredentials, details url.Values,token string) url.Values {
	form := make(url.Values, 12)
	form.Set("VERSION", string(Version124))
	form.Set("USER", credentials.Username)
	form.Set("PWD", credentials.Password)
	form.Set("SIGNATURE", credentials.Signature)
	form.Set("METHOD", "DoExpressCheckoutPayment")
	form.Set("TOKEN", token)
	form.Set("PAYMENTREQUEST_0_PAYMENTACTION", details.Get("PAYMENTREQUEST_0_PAYMENTACTION"))
	form.Set("PAYERID", details.Get("PAYERID"))
	form.Set("PAYMENTREQUEST_0_AMT", details.Get("PAYMENTREQUEST_0_AMT"))
	form.Set("PAYMENTREQUEST_0_ITEMAMT", details.Get("PAYMENTREQUEST_0_ITEMAMT"))
	form.Set("PAYMENTREQUEST_0_SHIPPINGAMT", details.Get("PAYMENTREQUEST_0_SHIPPINGAMT"))
	form.Set("PAYMENTREQUEST_0_TAXAMT", details.Get("PAYMENTREQUEST_0_TAXAMT"))
	form.Set("PAYMENTREQUEST_0_CURRENCYCODE", details.Get("PAYMENTREQUEST_0_CURRENCYCODE"))
	form.Set("PAYMENTREQUEST_0_DESC", details.Get("PAYMENTREQUEST_0_DESC"))
	form.Set("L_PAYMENTREQUEST_0_NAME0", details.Get("L_PAYMENTREQUEST_0_NAME0"))
	form.Set("L_PAYMENTREQUEST_0_AMT0", details.Get("L_PAYMENTREQUEST_0_AMT0"))
	form.Set("L_PAYMENTREQUEST_0_NUMBER0", details.Get("L_PAYMENTREQUEST_0_NUMBER0"))
	form.Set("L_PAYMENTREQUEST_0_QTY0", details.Get("L_PAYMENTREQUEST_0_QTY0"))
	form.Set("L_PAYMENTREQUEST_0_NAME1", details.Get("L_PAYMENTREQUEST_0_NAME1"))
	form.Set("L_PAYMENTREQUEST_0_AMT1", details.Get("L_PAYMENTREQUEST_0_AMT1"))
	form.Set("L_PAYMENTREQUEST_0_NUMBER1", details.Get("L_PAYMENTREQUEST_0_NUMBER1"))
	form.Set("L_PAYMENTREQUEST_0_QTY1", details.Get("L_PAYMENTREQUEST_0_QTY1"))
	// for duplicate
	form.Set("MSGSUBID", details.Get("PAYMENTREQUEST_0_INVNUM"))
	return form
}