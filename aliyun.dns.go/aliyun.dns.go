package aliyun

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

type RAM struct {
	AccessKeyId     string
	AccessKeySecret string
}

type DDNS struct {
	Ram      RAM
	RegionId string
}

const TTL = 1800
const identIPv6 = "http://v6.ipv6-test.com/api/myip.php"

/*
 * @param domainName: the domain name to be set
 * @param recordType: the record type to be set
 * @param rr: the record type to be set
 * @param interval: the interval to set
 */
func (dns *DDNS) AutoSetDomain(domain, resourceRecord, recordType, interval string) error {
	if dura, err := time.ParseDuration(interval); err != nil {
		return err
	} else {
		go func() {
			for {
				dns.SetDomain2Local(domain, resourceRecord, recordType)
				time.Sleep(dura)
			}
		}()
		return nil
	}
}

func (dns *DDNS) SetDomain2Local(domain, rr, recordType string) {

	ipv6, err := dns.GetMyIpv6()
	if err != nil {
		fmt.Println(err)
		return
	}
	record, err := dns.GetDomainRecordInfo(domain, rr, recordType)
	if err != nil {
		if err.Error() != "record not found" {
			fmt.Println(err)
			return
		} else {
			err := dns.SetDomainRecord(domain, rr, recordType, ipv6)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	if record.Value != ipv6 {

		fmt.Printf("%s Domain %s record %s type %s value %s is not equal to %s, updating...\n", timeStr, domain, rr, recordType, record.Value, ipv6)
		record.Value = ipv6
		if _, err := dns.UpdateDomianRecordInfo(record); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Printf("%s Domain %s record %s type %s value %s is equal to %s, no need to update.\n", timeStr, domain, rr, recordType, record.Value, ipv6)
	}
}

func (dns *DDNS) GetDomainRecordInfo(domain, rr, recordType string) (*alidns.Record, error) {
	client, err := alidns.NewClientWithAccessKey(dns.RegionId, dns.Ram.AccessKeyId, dns.Ram.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"

	request.DomainName = domain
	request.RRKeyWord = rr
	request.TypeKeyWord = recordType

	if response, err := client.DescribeDomainRecords(request); err != nil {
		return nil, err
	} else if len(response.DomainRecords.Record) > 0 {

		return &response.DomainRecords.Record[0], nil
	} else {
		return &alidns.Record{
			Value:      "{NULL}",
			DomainName: domain,
			RR:         rr,
			Type:       recordType,
			TTL:        TTL,
		}, fmt.Errorf("record not found")
	}
}

func (dns *DDNS) UpdateDomianRecordInfo(record *alidns.Record) (*alidns.UpdateDomainRecordResponse, error) {
	client, err := alidns.NewClientWithAccessKey(dns.RegionId, dns.Ram.AccessKeyId, dns.Ram.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"

	request.RecordId = record.RecordId
	request.RR = record.RR
	request.Type = record.Type
	request.Value = record.Value
	request.TTL = requests.NewInteger(int(record.TTL))

	return client.UpdateDomainRecord(request)
}
func (dns *DDNS) SetDomainRecord(domain, rr, recordType, value string) error {
	client, err := alidns.NewClientWithAccessKey(dns.RegionId, dns.Ram.AccessKeyId, dns.Ram.AccessKeySecret)
	if err != nil {
		return err
	}
	request := alidns.CreateAddDomainRecordRequest()
	request.Scheme = "https"

	request.DomainName = domain
	request.RR = rr
	request.Type = recordType
	request.Value = value
	request.TTL = requests.NewInteger(TTL)
	response, err := client.AddDomainRecord(request)
	if err != nil {
		return err
	}
	if !response.IsSuccess() {
		return fmt.Errorf("%s", response.BaseResponse.GetHttpContentString())
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s Domain %s record %s type %s value %s is created.\n", timeStr, domain, rr, recordType, value)
	return nil
}
func (dns *DDNS) GetMyIpv6() (string, error) {
	response, err := http.Get(identIPv6)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if body, err := ioutil.ReadAll(response.Body); err != nil {
		return "", err
	} else {
		return string(body), nil
	}
}
