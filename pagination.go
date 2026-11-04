package truto

import (
	"context"
	"iter"
	"net/http"
)

// Page is one page of a list endpoint.
type Page[T any] struct {
	Items []T
	// NextCursor is empty on the last page.
	NextCursor string
}

// Iterator walks a cursor-paginated list endpoint. It is lazy: no request is
// made until the first call to Next, NextPage, All, First or Seq.
//
//	it := client.IntegratedAccounts.List(ctx)
//	for {
//		acct, err := it.Next()
//		if err == truto.Done {
//			break
//		}
//		if err != nil {
//			return err
//		}
//		fmt.Println(*acct.ID)
//	}
//
// Or with Go 1.23 range-over-func:
//
//	for acct, err := range client.IntegratedAccounts.List(ctx).Seq() { ... }
//
// An Iterator is not safe for concurrent use.
type Iterator[T any] struct {
	ctx    context.Context
	client *Client
	path   string
	opts   []Option

	buf    []T
	cursor string
	done   bool
	err    error
}

func newIterator[T any](ctx context.Context, c *Client, path string, opts []Option) *Iterator[T] {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Iterator[T]{ctx: ctx, client: c, path: path, opts: opts}
}

// listResponse accepts both "result" (what the API returns) and "results"
// (what parts of the published spec say).
type listResponse[T any] struct {
	Result     []T     `json:"result"`
	Results    []T     `json:"results"`
	NextCursor *string `json:"next_cursor"`
}

// NextPage fetches the next page. After the last page it returns Done.
func (it *Iterator[T]) NextPage() (*Page[T], error) {
	if it.err != nil {
		return nil, it.err
	}
	if it.done {
		return nil, Done
	}
	opts := it.opts
	if it.cursor != "" {
		opts = append(append([]Option(nil), it.opts...), WithQuery(Params{"next_cursor": it.cursor}))
	}
	resp, err := it.client.do(it.ctx, request{method: http.MethodGet, path: it.path}, opts)
	if err != nil {
		it.err = err
		return nil, err
	}
	var lr listResponse[T]
	if err := resp.decodeInto(&lr); err != nil {
		it.err = err
		return nil, err
	}
	items := lr.Result
	if items == nil {
		items = lr.Results
	}
	it.cursor = ""
	if lr.NextCursor != nil {
		it.cursor = *lr.NextCursor
	}
	if it.cursor == "" {
		it.done = true
	}
	return &Page[T]{Items: items, NextCursor: it.cursor}, nil
}

// Next returns the next item, fetching pages as needed. It returns Done after
// the last item.
func (it *Iterator[T]) Next() (*T, error) {
	for len(it.buf) == 0 {
		if it.done || it.err != nil {
			if it.err != nil {
				return nil, it.err
			}
			return nil, Done
		}
		page, err := it.NextPage()
		if err != nil {
			return nil, err
		}
		it.buf = page.Items
	}
	item := it.buf[0]
	it.buf = it.buf[1:]
	return &item, nil
}

// All fetches every remaining page and returns all items.
func (it *Iterator[T]) All() ([]T, error) {
	var out []T
	for {
		item, err := it.Next()
		if err == Done {
			return out, nil
		}
		if err != nil {
			return out, err
		}
		out = append(out, *item)
	}
}

// First returns the first item of the first page, or Done if the list is empty.
func (it *Iterator[T]) First() (*T, error) {
	return it.Next()
}

// Seq returns a range-over-func sequence of the remaining items. Iteration
// stops after the first error, which is yielded with a zero item.
func (it *Iterator[T]) Seq() iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for {
			item, err := it.Next()
			if err == Done {
				return
			}
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}
			if !yield(*item, nil) {
				return
			}
		}
	}
}
