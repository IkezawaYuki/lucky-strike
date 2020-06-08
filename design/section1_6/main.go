package main

type Document struct {
}

type Machine interface {
	Print(d Document)
	Fax(d Document)
	Scan(d Document)
}

type MultiFunctionPrinter struct {
}

func (m MultiFunctionPrinter) Print(d Document) {
	panic("implement me")
}

func (m MultiFunctionPrinter) Fax(d Document) {
	panic("implement me")
}

func (m MultiFunctionPrinter) Scan(d Document) {
	panic("implement me")
}

type OldFashionedPrinter struct {
}

func (oo OldFashionedPrinter) Print(d Document) {
	// ok
}

func (oo OldFashionedPrinter) Fax(d Document) {
	panic("implement me")
}

func (oo OldFashionedPrinter) Scan(d Document) {
	panic("implement me")
}

type Printer interface {
	Print(d Document)
}

type Scanner interface {
	Scan(d Document)
}

type MyPrinter struct {
}

func (m MyPrinter) Print(d Document) {

}

type Photocopier struct {
}

func (p Photocopier) Scan(d Document) {

}

func (p Photocopier) Print(d Document) {

}

type MultiFunctionDevice interface {
	Printer
	Scanner
}

type MultiFunctionMachine struct {
	printer Printer
	scanner Scanner
}

func (m MultiFunctionMachine) Print(d Document) {
	m.printer.Print(d)
}

func (m MultiFunctionMachine) Scan(d Document) {
	m.scanner.Scan(d)
}

func main() {

}
