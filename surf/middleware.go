package surf

type (
	clientMiddleware   func(*Client)
	requestMiddleware  func(*Request) error
	responseMiddleware func(*Response) error
)

func (c *Client) applyReqMW(req *Request) error {
	for _, m := range c.reqMW {
		if err := m(req); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) applyRespMW(resp *Response) error {
	for _, m := range c.respMW {
		if err := m(resp); err != nil {
			return err
		}
	}

	return nil
}

func (opt *Options) applyReqMW(req *Request) error {
	for _, m := range opt.reqMW {
		if err := m(req); err != nil {
			return err
		}
	}

	return nil
}
