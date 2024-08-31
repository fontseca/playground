package playground

import (
  "bytes"
  "testing"
)

func TestFormatText(t *testing.T) {
  tests := [...][2]string{
    {"", ""},
    {`



    


      




package playground

import (
  "bytes"
  "html"
  "io"
  "log/slog"
  "regexp"
  "strings"
)

var (
  xmlReTag              = regexp.MustCompile('<([/!]?)([^>]+?)(/?)>')        // regexp xmlReTag matches an XML tag
  xmlReComment          = regexp.MustCompile('(?s)(<!--)(.*?)(-->)')         // regexp xmlReComment matches an XML comment
  xmlReXMLInsideComment = regexp.MustCompile('<!--[^>]*?<[^>]+?[^-]*?-->')   // regexp xmlReXMLInsideComment matches XML tags inside comments
  xmlReBlanksAround     = regexp.MustCompile('\s*(<|/?>)\s*')                // regexp xmlReBlanksAround matches blanks around XML tags
  xmlReBlanksInsideTags = regexp.MustCompile('>([^<]*[\n\r\t]|{3,})[^<]*<')  // regexp xmlReBlanksInsideTags matches the content of an XML tags that contains repeated spaces or \r\n\t
  xmlReBlanks           = regexp.MustCompile('\s{2,}')                       // regexp xmlReBlanks matches the content of an XML tag that contains more than two spaces together
)

type xmlFormatterImpl struct{}

// formatXML implements basic formatting for XML input by applying indentation, normalizing whitespace, handling
// comments and XML tags inside comments.
func (xmlFormatterImpl) format(input []byte, output io.Writer, indent string) {
  if nil == input || len(input) < 2 {
    return
  }

  var needsUnescape bool

  // Comments might contain further XML code. In that's the case, we want to escape
  // that code to avoid formatting.
  out := xmlReXMLInsideComment.ReplaceAllFunc(input, func(comment []byte) []byte {
    needsUnescape = true
    submatches := xmlReComment.FindSubmatch(comment)
    b := bytes.Buffer{}
    b.Grow(len(comment))
    b.Write(submatches[1])                                  // <!--
    b.WriteString(html.EscapeString(string(submatches[2]))) // ... (which includes XML code)
    b.Write(submatches[3])                                  // -->
    return b.Bytes()
  })

  out = xmlReBlanksAround.ReplaceAll(out, []byte("$1"))
  out = xmlReBlanksInsideTags.ReplaceAllFunc(out, func(m []byte) []byte {
    return xmlReBlanks.ReplaceAll(m, []byte(" "))
  })

  out = xmlReTag.ReplaceAllFunc(out, tagReplacer(indent))

  if needsUnescape {
    // restore the original comment escaped content
    out = xmlReComment.ReplaceAllFunc(out, func(comment []byte) []byte {
      submatches := xmlReComment.FindSubmatch(comment)
      b := bytes.Buffer{}
      b.Grow(len(comment))
      b.Write(submatches[1])                                    // <!--
      b.WriteString(html.UnescapeString(string(submatches[2]))) // ... (which include XML code)
      b.Write(submatches[3])                                    // -->
      return b.Bytes()
    })
  }

  _, err := output.Write(out[1:])
  if nil != err {
    slog.Error(err.Error())
  }
}















`, `package playground

import (
  "bytes"
  "html"
  "io"
  "log/slog"
  "regexp"
  "strings"
)

var (
  xmlReTag              = regexp.MustCompile('<([/!]?)([^>]+?)(/?)>')        // regexp xmlReTag matches an XML tag
  xmlReComment          = regexp.MustCompile('(?s)(<!--)(.*?)(-->)')         // regexp xmlReComment matches an XML comment
  xmlReXMLInsideComment = regexp.MustCompile('<!--[^>]*?<[^>]+?[^-]*?-->')   // regexp xmlReXMLInsideComment matches XML tags inside comments
  xmlReBlanksAround     = regexp.MustCompile('\s*(<|/?>)\s*')                // regexp xmlReBlanksAround matches blanks around XML tags
  xmlReBlanksInsideTags = regexp.MustCompile('>([^<]*[\n\r\t]|{3,})[^<]*<')  // regexp xmlReBlanksInsideTags matches the content of an XML tags that contains repeated spaces or \r\n\t
  xmlReBlanks           = regexp.MustCompile('\s{2,}')                       // regexp xmlReBlanks matches the content of an XML tag that contains more than two spaces together
)

type xmlFormatterImpl struct{}

// formatXML implements basic formatting for XML input by applying indentation, normalizing whitespace, handling
// comments and XML tags inside comments.
func (xmlFormatterImpl) format(input []byte, output io.Writer, indent string) {
  if nil == input || len(input) < 2 {
    return
  }

  var needsUnescape bool

  // Comments might contain further XML code. In that's the case, we want to escape
  // that code to avoid formatting.
  out := xmlReXMLInsideComment.ReplaceAllFunc(input, func(comment []byte) []byte {
    needsUnescape = true
    submatches := xmlReComment.FindSubmatch(comment)
    b := bytes.Buffer{}
    b.Grow(len(comment))
    b.Write(submatches[1])                                  // <!--
    b.WriteString(html.EscapeString(string(submatches[2]))) // ... (which includes XML code)
    b.Write(submatches[3])                                  // -->
    return b.Bytes()
  })

  out = xmlReBlanksAround.ReplaceAll(out, []byte("$1"))
  out = xmlReBlanksInsideTags.ReplaceAllFunc(out, func(m []byte) []byte {
    return xmlReBlanks.ReplaceAll(m, []byte(" "))
  })

  out = xmlReTag.ReplaceAllFunc(out, tagReplacer(indent))

  if needsUnescape {
    // restore the original comment escaped content
    out = xmlReComment.ReplaceAllFunc(out, func(comment []byte) []byte {
      submatches := xmlReComment.FindSubmatch(comment)
      b := bytes.Buffer{}
      b.Grow(len(comment))
      b.Write(submatches[1])                                    // <!--
      b.WriteString(html.UnescapeString(string(submatches[2]))) // ... (which include XML code)
      b.Write(submatches[3])                                    // -->
      return b.Bytes()
    })
  }

  _, err := output.Write(out[1:])
  if nil != err {
    slog.Error(err.Error())
  }
}`},
    {
      `formatter := textFormatterImpl{}
  for _, test := range tests {
    buf := bytes.Buffer{}
    expected := test[1]
    formatter.format([]byte(test[0]), &buf, "  ")
    got := string(buf.Bytes())
    if got != expected {
      t.Errorf("\n"+
        "\nexpected:\n\n%s\n"+
        "\ngot:\n\n%s",
        expected, got)
    }
  }`,
      `formatter := textFormatterImpl{}
  for _, test := range tests {
    buf := bytes.Buffer{}
    expected := test[1]
    formatter.format([]byte(test[0]), &buf, "  ")
    got := string(buf.Bytes())
    if got != expected {
      t.Errorf("\n"+
        "\nexpected:\n\n%s\n"+
        "\ngot:\n\n%s",
        expected, got)
    }
  }`},
  }

  formatter := textFormatterImpl{}
  for _, test := range tests {
    buf := bytes.Buffer{}
    expected := test[1]
    formatter.format([]byte(test[0]), &buf, "  ")
    got := string(buf.Bytes())
    if got != expected {
      t.Errorf("\n"+
        "\nexpected:\n\n%s\n"+
        "\ngot:\n\n%s",
        expected, got)
    }
  }
}

func TestFormatXML(t *testing.T) {
  tests := [...][2]string{
    {"", ""},
    {
      `<?xml  version="1.0"
              encoding="UTF-8"?><俄语 լեզու="ռուսերեն">данные</俄语>`,
      `<?xml version="1.0" encoding="UTF-8"?>
<俄语 լեզու="ռուսերեն">данные</俄语>`},
    {
      `<?xml version="1.0" encoding="UTF-8" ?><xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema"></xs:schema>`,
      `<?xml version="1.0" encoding="UTF-8" ?>
<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema"></xs:schema>`,
    },
    {`<objects type="array">
<object><id type="integer">1</id><name>Karl Zboncak</name><email>viva@keebler.biz</email><avatar>https://robohash.org/sequiquiaporro.png?size=300x300&set=set1</avatar>
<created-at>2021-06-24T15:45:46.667Z</created-at><updated-at>2021-06-24T15:45:46.667Z</updated-at></object></objects>`,
      `<objects type="array">
  <object>
    <id type="integer">1</id>
    <name>Karl Zboncak</name>
    <email>viva@keebler.biz</email>
    <avatar>https://robohash.org/sequiquiaporro.png?size=300x300&set=set1</avatar>
    <created-at>2021-06-24T15:45:46.667Z</created-at>
    <updated-at>2021-06-24T15:45:46.667Z</updated-at>
  </object>
</objects>`},
    {
      `<root><item/></root>`,
      `<root>
  <item/>
</root>`,
    },
    {
      `<root><level1><level2><level3><level4></level4></level3></level2></level1></root>`,
      `<root>
  <level1>
    <level2>
      <level3>
        <level4></level4>
      </level3>
    </level2>
  </level1>
</root>`,
    },
    {
      `<root><empty>    </empty></root>`,
      `<root>
  <empty></empty>
</root>`,
    },
    {
      `<object>   a     </object>
<object>    b    <  /object>`,
      `<object>a</object>
<object>b</object>`,
    },
    {
      `<root><!-- comment --><child><!-- nested comment --></child></root>`,
      `<root>
  <!-- comment -->
  <child>
    <!-- nested comment -->
  </child>
</root>`,
    },
    {`<case>
  <title>Lorem ipsum

dolor sit amet ... </title>
 
  <!-- The type of the test case. -->   
  <type>Usability</type> 
 
  <!-- The priority of the test case, specified as a simple
  number. -->
  <priority>2</priority>
 
  <!-- The estimate of the test case, specified in seconds. -->
  <estimate>100</estimate>
 
  <!-- The milestone of the test case. Is looked up in the
       project you import the test cases to. -->
  <milestone>Beta 1</milestone>
 
  <!-- The references of the test case (a list of IDs, separated
       by comma or space). May be used to link the test case to external
       references (requirements or user stories, e.g.). -->
  <references>REQ-1, REQ-2</references>
 
  <custom>
    <!-- Custom fields go here -->
  </custom>
        <custom2>        <!-- Custom fields go here -->      </custom2>
</case>`,
      `<case>
  <title>Lorem ipsum dolor sit amet ...</title>
  <!-- The type of the test case. -->
  <type>Usability</type>
  <!-- The priority of the test case, specified as a simple
  number. -->
  <priority>2</priority>
  <!-- The estimate of the test case, specified in seconds. -->
  <estimate>100</estimate>
  <!-- The milestone of the test case. Is looked up in the
       project you import the test cases to. -->
  <milestone>Beta 1</milestone>
  <!-- The references of the test case (a list of IDs, separated
       by comma or space). May be used to link the test case to external
       references (requirements or user stories, e.g.). -->
  <references>REQ-1, REQ-2</references>
  <custom>
    <!-- Custom fields go here -->
  </custom>
  <custom2>
    <!-- Custom fields go here -->
  </custom2>
</case>`},
    {
      `<?xml version="1.0"   

encoding="UTF-8"?>
<!DOCTYPE        people_list     SYSTEM "example.dtd">
<people_list>  <person>    <name>Fred Bloggs</name>    <birthdate>27/11/2008</birthdate>
    <gender>Male</gender>  </person>
</people_list>`,
      `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE people_list SYSTEM "example.dtd">
<people_list>
  <person>
    <name>Fred Bloggs</name>
    <birthdate>27/11/2008</birthdate>
    <gender>Male</gender>
  </person>
</people_list>`,
    },
    {`
<!-- the xml snippet here
     MUST not be formatted

       <people_list>
  <person><name>Fred Bloggs</name>    <birthdate>27/11/2008</birthdate>
    <gender>Male</gender>  </person>
</people_list>

    -->
<people_list>  <person>    <name>Fred Bloggs</name>    <birthdate>27/11/2008</birthdate>
    <gender>Male</gender>  </person>
</people_list>`,
      `<!-- the xml snippet here
     MUST not be formatted

       <people_list>
  <person><name>Fred Bloggs</name>    <birthdate>27/11/2008</birthdate>
    <gender>Male</gender>  </person>
</people_list>

    -->
<people_list>
  <person>
    <name>Fred Bloggs</name>
    <birthdate>27/11/2008</birthdate>
    <gender>Male</gender>
  </person>
</people_list>`},
  }

  formatter := xmlFormatterImpl{}
  for _, test := range tests {
    buf := bytes.Buffer{}
    expected := test[1]
    formatter.format([]byte(test[0]), &buf, "  ")
    got := string(buf.Bytes())
    if got != expected {
      t.Errorf("\n"+
        "\nexpected:\n\n%s\n"+
        "\ngot:\n\n%s",
        expected, got)
    }
  }
}

func BenchmarkFormatXML(b *testing.B) {
  input := []byte(`
<?xml version="1.0" encoding="UTF-8"?>  <Invoice
xmlns="urn:oasis:names:specification:ubl:schema:xsd:Invoice-2"
xmlns:cac="urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2"
xmlns:cbc="urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2">  <cbc:UBLVersionID>2.1</cbc:UBLVersionID>  <cbc:ID>TOSL108</cbc:ID>  <cbc:IssueDate>2009-12-15</cbc:IssueDate>  <cbc:InvoiceTypeCode listID="UN/ECE 1001 Subset" listAgencyID="6">380</cbc:InvoiceTypeCode>  <cbc:Note languageID="en">Ordered in our booth at the convention.</cbc:Note>  <cbc:TaxPointDate>2009-11-30</cbc:TaxPointDate>  <cbc:DocumentCurrencyCode listID="ISO 4217 Alpha" listAgencyID="6">EUR</cbc:DocumentCurrencyCode>  <cbc:AccountingCost>Project cost code 123</cbc:AccountingCost>  <cac:InvoicePeriod>  <cbc:StartDate>2009-11-01</cbc:StartDate> 
<cbc:EndDate>2009-11-30</cbc:EndDate>  </cac:InvoicePeriod>  <cac:OrderReference>  <cbc:ID>123</cbc:ID>  </cac:OrderReference>  <cac:ContractDocumentReference>  <cbc:ID>Contract321</cbc:ID>  <cbc:DocumentType>Framework agreement</cbc:DocumentType>  </cac:ContractDocumentReference> 
<cac:AdditionalDocumentReference>  <cbc:ID>Doc1</cbc:ID>  <cbc:DocumentType>Timesheet</cbc:DocumentType>  <cac:Attachment>  <cac:ExternalReference>  <cbc:URI>http://www.suppliersite.eu/sheet001.html</cbc:URI>  </cac:ExternalReference>  </cac:Attachment>  </cac:AdditionalDocumentReference>  <cac:AdditionalDocumentReference>  <cbc:ID>Doc2</cbc:ID>  <cbc:DocumentType>Drawing</cbc:DocumentType>  <cac:Attachment>  <cbc:EmbeddedDocumentBinaryObject
mimeCode="application/pdf">UjBsR09EbGhjZ0dTQUxNQUFBUUNBRU1tQ1p0dU1GUXhEUzhi</cbc:EmbeddedDocumentBinaryObject>  </cac:Attachment>  </cac:AdditionalDocumentReference>  <cac:AccountingSupplierParty>  <cac:Party>  <cbc:EndpointID schemeID="GLN" schemeAgencyID="9">1234567890123</cbc:EndpointID>  <cac:PartyIdentification>  <cbc:ID schemeID="ZZZ">Supp123</cbc:ID>  </cac:PartyIdentification>  <cac:PartyName>  <cbc:Name>Salescompany ltd.</cbc:Name>  </cac:PartyName>  <cac:PostalAddress>  <cbc:ID schemeID="GLN" schemeAgencyID="9">1231412341324</cbc:ID>  <cbc:Postbox>5467</cbc:Postbox>  <cbc:StreetName>Main street</cbc:StreetName>  <cbc:AdditionalStreetName>Suite 123</cbc:AdditionalStreetName>  <cbc:BuildingNumber>1</cbc:BuildingNumber>  <cbc:Department>Revenue department</cbc:Department>  <cbc:CityName>Big city</cbc:CityName>  <cbc:PostalZone>54321</cbc:PostalZone>  <cbc:CountrySubentityCode>RegionA</cbc:CountrySubentityCode>  <cac:Country>  <cbc:IdentificationCode listID="ISO3166-1" listAgencyID="6">DK</cbc:IdentificationCode>  </cac:Country>  </cac:PostalAddress>  <cac:PartyTaxScheme>  <cbc:CompanyID schemeID="DKVAT" schemeAgencyID="ZZZ">DK12345</cbc:CompanyID>  <cac:TaxScheme>  <cbc:ID schemeID="UN/ECE 5153" schemeAgencyID="6">VAT</cbc:ID>  </cac:TaxScheme>  </cac:PartyTaxScheme>  <cac:PartyLegalEntity>  <cbc:RegistrationName>The Sellercompany Incorporated</cbc:RegistrationName>  <cbc:CompanyID schemeID="CVR" schemeAgencyID="ZZZ">5402697509</cbc:CompanyID>  <cac:RegistrationAddress>  <cbc:CityName>Big city</cbc:CityName>  <cbc:CountrySubentity>RegionA</cbc:CountrySubentity>  <cac:Country>  <cbc:IdentificationCode>DK</cbc:IdentificationCode>  </cac:Country>  </cac:RegistrationAddress>  </cac:PartyLegalEntity>  <cac:Contact>  <cbc:Telephone>4621230</cbc:Telephone>  <cbc:Telefax>4621231</cbc:Telefax>  <cbc:ElectronicMail>antonio@salescompany.dk</cbc:ElectronicMail>  </cac:Contact>  <cac:Person>  <cbc:FirstName>Antonio</cbc:FirstName>  <cbc:FamilyName>M</cbc:FamilyName>  <cbc:MiddleName>Salemacher</cbc:MiddleName>  <cbc:JobTitle>Sales manager</cbc:JobTitle>  </cac:Person>  </cac:Party>  </cac:AccountingSupplierParty>  <cac:AccountingCustomerParty>  <cac:Party>  <cbc:EndpointID schemeID="GLN" schemeAgencyID="9">1234567987654</cbc:EndpointID>  <cac:PartyIdentification>  <cbc:ID schemeID="ZZZ">345KS5324</cbc:ID>  </cac:PartyIdentification>  <cac:PartyName>  <cbc:Name>Buyercompany ltd</cbc:Name>  </cac:PartyName>  <cac:PostalAddress>  <cbc:ID schemeID="GLN" schemeAgencyID="9">1238764941386</cbc:ID>  <cbc:Postbox>123</cbc:Postbox>  <cbc:StreetName>Anystreet</cbc:StreetName>  <cbc:AdditionalStreetName>Back door</cbc:AdditionalStreetName>  <cbc:BuildingNumber>8</cbc:BuildingNumber>  <cbc:Department>Accounting department</cbc:Department>  <cbc:CityName>Anytown</cbc:CityName>  <cbc:PostalZone>101</cbc:PostalZone>  <cbc:CountrySubentity>RegionB</cbc:CountrySubentity>  <cac:Country>  <cbc:IdentificationCode listID="ISO3166-1" listAgencyID="6">BE</cbc:IdentificationCode>  </cac:Country>  </cac:PostalAddress>  <cac:PartyTaxScheme>  <cbc:CompanyID schemeID="BEVAT" schemeAgencyID="ZZZ">BE54321</cbc:CompanyID>  <cac:TaxScheme>  <cbc:ID schemeID="UN/ECE 5153" schemeAgencyID="6">VAT</cbc:ID>  </cac:TaxScheme>  </cac:PartyTaxScheme>  <cac:PartyLegalEntity>  <cbc:RegistrationName>The buyercompany inc.</cbc:RegistrationName>  <cbc:CompanyID schemeAgencyID="ZZZ" schemeID="ZZZ">5645342123</cbc:CompanyID>  <cac:RegistrationAddress>  <cbc:CityName>Mainplace</cbc:CityName>  <cbc:CountrySubentity>RegionB</cbc:CountrySubentity>  <cac:Country>  <cbc:IdentificationCode>BE</cbc:IdentificationCode>  </cac:Country>  </cac:RegistrationAddress>  </cac:PartyLegalEntity>  <cac:Contact>  <cbc:Telephone>5121230</cbc:Telephone>  <cbc:Telefax>5121231</cbc:Telefax>  <cbc:ElectronicMail>john@buyercompany.eu</cbc:ElectronicMail>  </cac:Contact>  <cac:Person>  <cbc:FirstName>John</cbc:FirstName>  <cbc:FamilyName>X</cbc:FamilyName>  <cbc:MiddleName>Doe</cbc:MiddleName>  <cbc:JobTitle>Purchasing manager</cbc:JobTitle>  </cac:Person>  </cac:Party>  </cac:AccountingCustomerParty>  <cac:PayeeParty>  <cac:PartyIdentification>  <cbc:ID schemeID="GLN" schemeAgencyID="9">098740918237</cbc:ID>  </cac:PartyIdentification>  <cac:PartyName>  <cbc:Name>Ebeneser Scrooge Inc.</cbc:Name>  </cac:PartyName>  <cac:PartyLegalEntity>  <cbc:CompanyID schemeID="UK:CH" schemeAgencyID="ZZZ">6411982340</cbc:CompanyID>  </cac:PartyLegalEntity>  </cac:PayeeParty>  <cac:Delivery>  <cbc:ActualDeliveryDate>2009-12-15</cbc:ActualDeliveryDate>  <cac:DeliveryLocation>  <cbc:ID schemeID="GLN" schemeAgencyID="9">6754238987648</cbc:ID>  <cac:Address>  <cbc:StreetName>Deliverystreet</cbc:StreetName>  <cbc:AdditionalStreetName>Side door</cbc:AdditionalStreetName>  <cbc:BuildingNumber>12</cbc:BuildingNumber>  <cbc:CityName>DeliveryCity</cbc:CityName>  <cbc:PostalZone>523427</cbc:PostalZone>  <cbc:CountrySubentity>RegionC</cbc:CountrySubentity>  <cac:Country>  <cbc:IdentificationCode>BE</cbc:IdentificationCode>  </cac:Country>  </cac:Address>  </cac:DeliveryLocation>  </cac:Delivery>  <cac:PaymentMeans>  <cbc:PaymentMeansCode listID="UN/ECE 4461">31</cbc:PaymentMeansCode>  <cbc:PaymentDueDate>2009-12-31</cbc:PaymentDueDate>  <cbc:PaymentChannelCode>IBAN</cbc:PaymentChannelCode>  <cbc:PaymentID>Payref1</cbc:PaymentID>  <cac:PayeeFinancialAccount>  <cbc:ID>DK1212341234123412</cbc:ID>  <cac:FinancialInstitutionBranch>  <cac:FinancialInstitution>  <cbc:ID>DKDKABCD</cbc:ID>  </cac:FinancialInstitution>  </cac:FinancialInstitutionBranch>  </cac:PayeeFinancialAccount>  </cac:PaymentMeans>  <cac:PaymentTerms>  <cbc:Note>Penalty percentage 10% from due date</cbc:Note>  </cac:PaymentTerms>  <cac:AllowanceCharge>  <cbc:ChargeIndicator>true</cbc:ChargeIndicator>  <cbc:AllowanceChargeReason>Packing cost</cbc:AllowanceChargeReason>  <cbc:Amount currencyID="EUR">100</cbc:Amount>  </cac:AllowanceCharge>  <cac:AllowanceCharge>  <cbc:ChargeIndicator>false</cbc:ChargeIndicator>  <cbc:AllowanceChargeReason>Promotion discount</cbc:AllowanceChargeReason>  <cbc:Amount currencyID="EUR">100</cbc:Amount>  </cac:AllowanceCharge>  <cac:TaxTotal>  <cbc:TaxAmount currencyID="EUR">292.20</cbc:TaxAmount>  <cac:TaxSubtotal>  <cbc:TaxableAmount currencyID="EUR">1460.5</cbc:TaxableAmount>  <cbc:TaxAmount currencyID="EUR">292.1</cbc:TaxAmount>  <cac:TaxCategory>  <cbc:ID schemeID="UN/ECE 5305" schemeAgencyID="6">S</cbc:ID>  <cbc:Percent>20</cbc:Percent>  <cac:TaxScheme>  <cbc:ID schemeID="UN/ECE 5153" schemeAgencyID="6">VAT</cbc:ID>  </cac:TaxScheme>  </cac:TaxCategory>  </cac:TaxSubtotal>  <cac:TaxSubtotal>  <cbc:TaxableAmount currencyID="EUR">1</cbc:TaxableAmount>  <cbc:TaxAmount currencyID="EUR">0.1</cbc:TaxAmount>  <cac:TaxCategory>  <cbc:ID schemeID="UN/ECE 5305" schemeAgencyID="6">AA</cbc:ID>  <cbc:Percent>10</cbc:Percent>  <cac:TaxScheme>  <cbc:ID schemeID="UN/ECE 5153" schemeAgencyID="6">VAT</cbc:ID>  </cac:TaxScheme>  </cac:TaxCategory>  </cac:TaxSubtotal>  <cac:TaxSubtotal>  <cbc:TaxableAmount currencyID="EUR">-25</cbc:TaxableAmount>  <cbc:TaxAmount currencyID="EUR">0</cbc:TaxAmount>  <cac:TaxCategory>  <cbc:ID schemeID="UN/ECE 5305" schemeAgencyID="6">E</cbc:ID>  <cbc:Percent>0</cbc:Percent>  <cbc:TaxExemptionReasonCode listID="CWA 15577" listAgencyID="ZZZ">AAM</cbc:TaxExemptionReasonCode>  <cbc:TaxExemptionReason>Exempt New Means of Transport</cbc:TaxExemptionReason>  <cac:TaxScheme>  <cbc:ID schemeID="UN/ECE 5153" schemeAgencyID="6">VAT</cbc:ID>  </cac:TaxScheme>  </cac:TaxCategory>  </cac:TaxSubtotal>  </cac:TaxTotal>  <cac:LegalMonetaryTotal>  <cbc:LineExtensionAmount currencyID="EUR">1436.5</cbc:LineExtensionAmount>  <cbc:TaxExclusiveAmount currencyID="EUR">1436.5</cbc:TaxExclusiveAmount>  <cbc:TaxInclusiveAmount currencyID="EUR">1729</cbc:TaxInclusiveAmount>  <cbc:AllowanceTotalAmount currencyID="EUR">100</cbc:AllowanceTotalAmount>  <cbc:ChargeTotalAmount currencyID="EUR">100</cbc:ChargeTotalAmount>  <cbc:PrepaidAmount currencyID="EUR">1000</cbc:PrepaidAmount>  <cbc:PayableRoundingAmount currencyID="EUR">0.30</cbc:PayableRoundingAmount>  <cbc:PayableAmount currencyID="EUR">729</cbc:PayableAmount>  </cac:LegalMonetaryTotal>  <cac:InvoiceLine>  <cbc:ID>1</cbc:ID>  <cbc:Note>Scratch on box</cbc:Note>  <cbc:InvoicedQuantity unitCode="C62">1</cbc:InvoicedQuantity>  <cbc:LineExtensionAmount currencyID="EUR">1273</cbc:LineExtensionAmount>  <cbc:AccountingCost>BookingCode001</cbc:AccountingCost>  <cac:OrderLineReference>  <cbc:LineID>1</cbc:LineID>  </cac:OrderLineReference>  <cac:AllowanceCharge>  <cbc:ChargeIndicator>false</cbc:ChargeIndicator>  <cbc:AllowanceChargeReason>Damage</cbc:AllowanceChargeReason>  <cbc:Amount currencyID="EUR">12</cbc:Amount>  </cac:AllowanceCharge>  <cac:AllowanceCharge>  <cbc:ChargeIndicator>true</cbc:ChargeIndicator>  <cbc:AllowanceChargeReason>Testing</cbc:AllowanceChargeReason>  <cbc:Amount currencyID="EUR">10</cbc:Amount>  </cac:AllowanceCharge>  <cac:TaxTotal>  <cbc:TaxAmount currencyID="EUR">254.6</cbc:TaxAmount>  </cac:TaxTotal>  <cac:Item>  <cbc:Description languageID="EN">Processor: Intel Core 2 Duo SU9400 LV (1.4GHz). RAM: 3MB. Screen 1440x900</cbc:Description>  <cbc:Name>Labtop computer</cbc:Name>  <cac:SellersItemIdentification>  <cbc:ID>JB007</cbc:ID>  </cac:SellersItemIdentification>  <cac:StandardItemIdentification>  <cbc:ID schemeID="GTIN" schemeAgencyID="9">1234567890124</cbc:ID>  </cac:StandardItemIdentification>  <cac:CommodityClassification>  <cbc:ItemClassificationCode listAgencyID="113" listID="UNSPSC">12344321</cbc:ItemClassificationCode>  </cac:CommodityClassification>  <cac:CommodityClassification>  <cbc:ItemClassificationCode listAgencyID="2" listID="CPV">65434568</cbc:ItemClassificationCode>  </cac:CommodityClassification>  <ClassifiedTaxCategory xmlns="urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2">  <cbc:ID schemeID="UN/ECE 5305" schemeAgencyID="6">S</cbc:ID>  <cbc:Percent>20</cbc:Percent>  <TaxScheme>  <cbc:ID schemeID="UN/ECE 5153" schemeAgencyID="6">VAT</cbc:ID>  </TaxScheme>  </ClassifiedTaxCategory>  <cac:AdditionalItemProperty>  <cbc:Name>Color</cbc:Name>  <cbc:Value>black</cbc:Value>  </cac:AdditionalItemProperty>  </cac:Item>  <cac:Price>  <cbc:PriceAmount currencyID="EUR">1273</cbc:PriceAmount>  <cbc:BaseQuantity unitCode="C62">1</cbc:BaseQuantity>  <cac:AllowanceCharge>  <cbc:ChargeIndicator>false</cbc:ChargeIndicator>  <cbc:AllowanceChargeReason>Contract</cbc:AllowanceChargeReason>  <cbc:MultiplierFactorNumeric>0.15</cbc:MultiplierFactorNumeric>  <cbc:Amount currencyID="EUR">225</cbc:Amount>  <cbc:BaseAmount currencyID="EUR">1500</cbc:BaseAmount>  </cac:AllowanceCharge>  </cac:Price>  </cac:InvoiceLine>  <cac:InvoiceLine>  <cbc:ID>2</cbc:ID>  <cbc:Note>Cover is slightly damaged.</cbc:Note>  <cbc:InvoicedQuantity unitCode="C62">-1</cbc:InvoicedQuantity>  <cbc:LineExtensionAmount currencyID="EUR">-3.96</cbc:LineExtensionAmount>  <cac:OrderLineReference>  <cbc:LineID>5</cbc:LineID>  </cac:OrderLineReference>  <cac:TaxTotal>  <cbc:TaxAmount currencyID="EUR">-0.396</cbc:TaxAmount>  </cac:TaxTotal>  <cac:Item>  <cbc:Name>Returned "Advanced computing" book</cbc:Name>  <cac:SellersItemIdentification>  <cbc:ID>JB008</cbc:ID>  </cac:SellersItemIdentification>  <cac:StandardItemIdentification>  <cbc:ID schemeID="GTIN" schemeAgencyID="9">1234567890125</cbc:ID>  </cac:StandardItemIdentification>  <cac:CommodityClassification>  <cbc:ItemClassificationCode listAgencyID="113" listID="UNSPSC">32344324</cbc:ItemClassificationCode>  </cac:CommodityClassification>  <cac:CommodityClassification>  <cbc:ItemClassificationCode listAgencyID="2" listID="CPV">65434567</cbc:ItemClassificationCode>  </cac:CommodityClassification>  <ClassifiedTaxCategory xmlns="urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2">  <cbc:ID schemeID="UN/ECE 5305" schemeAgencyID="6">AA</cbc:ID>  <cbc:Percent>10</cbc:Percent>  <TaxScheme>  <cbc:ID schemeID="UN/ECE 5153" schemeAgencyID="6">VAT</cbc:ID>  </TaxScheme>  </ClassifiedTaxCategory>  </cac:Item>  <cac:Price>  <cbc:PriceAmount currencyID="EUR">3.96</cbc:PriceAmount>  <cbc:BaseQuantity unitCode="C62">1</cbc:BaseQuantity>  </cac:Price>  </cac:InvoiceLine>  <cac:InvoiceLine>  <cbc:ID>3</cbc:ID>  <cbc:InvoicedQuantity unitCode="C62">2</cbc:InvoicedQuantity>  <cbc:LineExtensionAmount currencyID="EUR">4.96</cbc:LineExtensionAmount>  <cac:OrderLineReference>  <cbc:LineID>3</cbc:LineID>  </cac:OrderLineReference>  <cac:TaxTotal>  <cbc:TaxAmount currencyID="EUR">0.496</cbc:TaxAmount>  </cac:TaxTotal>  <cac:Item>  <cbc:Name>"Computing for dummies" book</cbc:Name>  <cac:SellersItemIdentification>  <cbc:ID>JB009</cbc:ID>  </cac:SellersItemIdentification>  <cac:StandardItemIdentification>  <cbc:ID schemeID="GTIN" schemeAgencyID="9">1234567890126</cbc:ID>  </cac:StandardItemIdentification>  <cac:CommodityClassification>  <cbc:ItemClassificationCode listAgencyID="113" listID="UNSPSC">32344324</cbc:ItemClassificationCode>  </cac:CommodityClassification>  <cac:CommodityClassification>  <cbc:ItemClassificationCode listAgencyID="2" listID="CPV">65434566</cbc:ItemClassificationCode>  </cac:CommodityClassification>  <ClassifiedTaxCategory xmlns="urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2">  <cbc:ID schemeID="UN/ECE 5305" schemeAgencyID="6">AA</cbc:ID>  <cbc:Percent>10</cbc:Percent>  <TaxScheme>  <cbc:ID schemeID="UN/ECE 5153" schemeAgencyID="6">VAT</cbc:ID>  </TaxScheme>  </ClassifiedTaxCategory>  </cac:Item>  <cac:Price>  <cbc:PriceAmount currencyID="EUR">2.48</cbc:PriceAmount>  <cbc:BaseQuantity unitCode="C62">1</cbc:BaseQuantity>  <cac:AllowanceCharge>  <cbc:ChargeIndicator>false</cbc:ChargeIndicator>  <cbc:AllowanceChargeReason>Contract</cbc:AllowanceChargeReason>  <cbc:MultiplierFactorNumeric>0.1</cbc:MultiplierFactorNumeric>  <cbc:Amount currencyID="EUR">0.275</cbc:Amount>  <cbc:BaseAmount currencyID="EUR">2.75</cbc:BaseAmount>  </cac:AllowanceCharge>  </cac:Price>  </cac:InvoiceLine>  <cac:InvoiceLine>  <cbc:ID>4</cbc:ID>  <cbc:InvoicedQuantity unitCode="C62">-1</cbc:InvoicedQuantity>  <cbc:LineExtensionAmount currencyID="EUR">-25</cbc:LineExtensionAmount>  <cac:OrderLineReference>  <cbc:LineID>2</cbc:LineID>  </cac:OrderLineReference>  <cac:TaxTotal>  <cbc:TaxAmount currencyID="EUR">0</cbc:TaxAmount>  </cac:TaxTotal>  <cac:Item>  <cbc:Name>Returned IBM 5150 desktop</cbc:Name>  <cac:SellersItemIdentification>  <cbc:ID>JB010</cbc:ID>  </cac:SellersItemIdentification>  <cac:StandardItemIdentification>  <cbc:ID schemeID="GTIN" schemeAgencyID="9">1234567890127</cbc:ID>  </cac:StandardItemIdentification>  <cac:CommodityClassification>  <cbc:ItemClassificationCode listAgencyID="113" listID="UNSPSC">12344322</cbc:ItemClassificationCode>  </cac:CommodityClassification>  <cac:CommodityClassification>  <cbc:ItemClassificationCode listAgencyID="2" listID="CPV">65434565</cbc:ItemClassificationCode>  </cac:CommodityClassification>  <ClassifiedTaxCategory xmlns="urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2">  <cbc:ID schemeID="UN/ECE 5305" schemeAgencyID="6">E</cbc:ID>  <cbc:Percent>0</cbc:Percent>  <TaxScheme>  <cbc:ID schemeID="UN/ECE 5153" schemeAgencyID="6">VAT</cbc:ID>  </TaxScheme>  </ClassifiedTaxCategory>  </cac:Item>  <cac:Price>  <cbc:PriceAmount currencyID="EUR">25</cbc:PriceAmount>  <cbc:BaseQuantity unitCode="C62">1</cbc:BaseQuantity>  </cac:Price>  </cac:InvoiceLine>  <cac:InvoiceLine>  <cbc:ID>5</cbc:ID>  <cbc:InvoicedQuantity unitCode="C62">250</cbc:InvoicedQuantity>  <cbc:LineExtensionAmount currencyID="EUR">187.5</cbc:LineExtensionAmount>  <cbc:AccountingCost>BookingCode002</cbc:AccountingCost>  <cac:OrderLineReference>  <cbc:LineID>4</cbc:LineID>  </cac:OrderLineReference>  <cac:TaxTotal>  <cbc:TaxAmount currencyID="EUR">37.5</cbc:TaxAmount>  </cac:TaxTotal>  <cac:Item>  <cbc:Name>Network cable</cbc:Name>  <cac:SellersItemIdentification>  <cbc:ID>JB011</cbc:ID>  </cac:SellersItemIdentification>  <cac:StandardItemIdentification>  <cbc:ID schemeID="GTIN" schemeAgencyID="9">1234567890128</cbc:ID>  </cac:StandardItemIdentification>  <cac:CommodityClassification>  <cbc:ItemClassificationCode listAgencyID="113" listID="UNSPSC">12344325</cbc:ItemClassificationCode>  </cac:CommodityClassification>  <cac:CommodityClassification>  <cbc:ItemClassificationCode listAgencyID="2" listID="CPV">65434564</cbc:ItemClassificationCode>  </cac:CommodityClassification>  <ClassifiedTaxCategory xmlns="urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2">  <cbc:ID schemeID="UN/ECE 5305" schemeAgencyID="6">S</cbc:ID>  <cbc:Percent>20</cbc:Percent>  <TaxScheme>  <cbc:ID schemeID="UN/ECE 5153" schemeAgencyID="6">VAT</cbc:ID>  </TaxScheme>  </ClassifiedTaxCategory>  <cac:AdditionalItemProperty>  <cbc:Name>Type</cbc:Name>  <cbc:Value>Cat5</cbc:Value>  </cac:AdditionalItemProperty>  </cac:Item>  <cac:Price>  <cbc:PriceAmount currencyID="EUR">0.75</cbc:PriceAmount>  <cbc:BaseQuantity unitCode="C62">1</cbc:BaseQuantity>  </cac:Price>  </cac:InvoiceLine>  </Invoice>
`)
  formatter := xmlFormatterImpl{}
  buf := &bytes.Buffer{}
  buf.Grow(512)
  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    formatter.format(input, buf, "  ")
  }
}
