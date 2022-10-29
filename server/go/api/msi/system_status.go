//go:build windows && amd64

package msi

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	msiModelsPb "github.com/s77rt/rdpcloud/proto/go/models/msi"
	"github.com/s77rt/rdpcloud/server/go/internal/encode"
	"github.com/s77rt/rdpcloud/server/go/internal/win"
	msiInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/msi"
)

func getProductsGuids(wg *sync.WaitGroup, productsGuids *[]string, err *error) {
	defer wg.Done()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var iProductIndex uint32 = 0

	for {
		lpProductBuf := make([]uint16, 38+1) // expected GUID format is "{xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}" of length 38, +1 for the terminating null character

		ret, _, _ := msiInternalApi.MsiEnumProductsW(
			iProductIndex,
			&lpProductBuf[0],
		)
		if ret != win.ERROR_SUCCESS {
			if ret != win.ERROR_NO_MORE_ITEMS {
				*err = status.Errorf(codes.Unknown, "Failed to get products guids (error: 0x%x)", ret)
				return
			}
			break
		}

		*productsGuids = append(*productsGuids, encode.UTF16ToString(lpProductBuf)[1:37])

		iProductIndex += 1
	}

	return
}

func getProductInfo(wg *sync.WaitGroup, productGuid string, productAttribute string, productAttributeValue *string, err *error) {
	defer wg.Done()

	szProduct, rerr := encode.UTF16PtrFromString(fmt.Sprintf("{%s}", productGuid))
	if rerr != nil {
		*err = status.Errorf(codes.InvalidArgument, "Unable to encode product guid to UTF16")
		return
	}

	szAttribute, rerr := encode.UTF16PtrFromString(productAttribute)
	if rerr != nil {
		*err = status.Errorf(codes.InvalidArgument, "Unable to encode product attribute to UTF16")
		return
	}

	var pcchValueBuf uint32

	ret, _, _ := msiInternalApi.MsiGetProductInfoW(
		szProduct,
		szAttribute,
		nil,
		&pcchValueBuf,
	)

	if ret != win.ERROR_SUCCESS && ret != win.ERROR_MORE_DATA {
		if ret == win.ERROR_UNKNOWN_PROPERTY {
			// The property does not exist. Nothing to do (early exit)
			return
		}
		*err = status.Errorf(codes.Unknown, "Failed to get product info (1) (error: 0x%x)", ret)
		return
	}

	lpValueBuf := make([]uint16, pcchValueBuf+1)
	pcchValueBuf = uint32(len(lpValueBuf))

	ret, _, _ = msiInternalApi.MsiGetProductInfoW(
		szProduct,
		szAttribute,
		&lpValueBuf[0],
		&pcchValueBuf,
	)

	if ret != win.ERROR_SUCCESS {
		if ret == win.ERROR_UNKNOWN_PROPERTY {
			// The property does not exist. Nothing to do (early exit)
			return
		}
		*err = status.Errorf(codes.Unknown, "Failed to get product info (2) (error: 0x%x)", ret)
		return
	}

	*productAttributeValue = encode.UTF16ToString(lpValueBuf)

	return
}

func GetProducts() ([]*msiModelsPb.Product, error) {
	var products []*msiModelsPb.Product

	var (
		productsGuids []string
		err           error
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go getProductsGuids(&wg, &productsGuids, &err)
	wg.Wait()

	if err != nil {
		return nil, err
	}

	for _, productGuid := range productsGuids {
		product := &msiModelsPb.Product{
			Guid: productGuid,
		}

		var (
			productName    string
			productNameErr error

			productVersion    string
			productVersionErr error

			productPublisher    string
			productPublisherErr error

			productInstallDate    string
			productInstallDateErr error
		)

		wg.Add(1)
		go getProductInfo(&wg, productGuid, "ProductName", &productName, &productNameErr)
		wg.Add(1)
		go getProductInfo(&wg, productGuid, "Version", &productVersion, &productVersionErr)
		wg.Add(1)
		go getProductInfo(&wg, productGuid, "Publisher", &productPublisher, &productPublisherErr)
		wg.Add(1)
		go getProductInfo(&wg, productGuid, "InstallDate", &productInstallDate, &productInstallDateErr)

		wg.Wait()

		if productNameErr != nil {
			log.Printf("An error occurred while gathering product info (%s) (ProductName): %v", productGuid, productNameErr)
		} else {
			product.Name = productName
		}

		if productVersionErr != nil {
			log.Printf("An error occurred while gathering product info (%s) (Version): %v", productGuid, productVersionErr)
		} else {
			product.Version = productVersion
		}

		if productPublisherErr != nil {
			log.Printf("An error occurred while gathering product info (%s) (Publisher): %v", productGuid, productPublisherErr)
		} else {
			product.Publisher = productPublisher
		}

		if productInstallDateErr != nil {
			log.Printf("An error occurred while gathering product info (%s) (InstallDate): %v", productGuid, productInstallDateErr)
		} else {
			if productInstallDate != "" {
				t, err := time.Parse("20060102", productInstallDate)
				if err != nil {
					log.Printf("An error occurred while gathering product info (%s) (InstallDate) (Parse): %v", productGuid, err)
				} else {
					product.InstallDate = timestamppb.New(t)
				}
			}
		}

		products = append(products, product)
	}

	return products, nil
}
